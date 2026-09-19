package ws

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopa/internal/core/domain"
	"gopa/internal/services"
	"gopa/pkg/constants"
	"gopa/pkg/utils"
	"gopa/pkg/utils/response"
)

type mockPomodoroStore struct {
	state *domain.PomodoroState
}

func (m *mockPomodoroStore) Get(_ context.Context, _ uuid.UUID) (domain.PomodoroState, error) {
	if m.state == nil {
		return domain.PomodoroState{}, domain.ErrNotFound
	}
	return *m.state, nil
}
func (m *mockPomodoroStore) Set(_ context.Context, _ uuid.UUID, s domain.PomodoroState) error {
	m.state = &s
	return nil
}
func (m *mockPomodoroStore) Delete(_ context.Context, _ uuid.UUID) error {
	m.state = nil
	return nil
}

type mockPomodoroHistoryRepo struct{}

func (mockPomodoroHistoryRepo) Create(_ context.Context, _ domain.PomodoroHistory) error {
	return nil
}
func (mockPomodoroHistoryRepo) List(_ context.Context, _ uuid.UUID, _ int, _ string) ([]domain.PomodoroHistory, error) {
	return nil, nil
}

func testAuthMiddleware(tokens *utils.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string
		header := c.GetHeader("Authorization")
		if header != "" {
			parts := strings.SplitN(header, " ", 2)
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				tokenStr = strings.TrimSpace(parts[1])
			}
		}
		if tokenStr == "" {
			tokenStr = strings.TrimSpace(c.Query("token"))
		}
		if tokenStr == "" {
			tokenStr = strings.TrimSpace(c.Query("access_token"))
		}
		if tokenStr == "" {
			response.Fail(c, domain.ErrUnauthorized)
			c.Abort()
			return
		}
		claims, err := tokens.ParseAccessToken(tokenStr)
		if err != nil {
			response.Fail(c, domain.ErrUnauthorized)
			c.Abort()
			return
		}
		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			response.Fail(c, domain.ErrUnauthorized)
			c.Abort()
			return
		}
		c.Set(constants.IdentityKey, userID)
		c.Set("user_id", userID)
		c.Next()
	}
}

func TestPomodoroWSGateway_OriginCheck(t *testing.T) {
	gateway := NewPomodoroWSGateway(nil, nil, slog.Default(), "http://localhost:5173")

	reqEmpty, _ := http.NewRequest(http.MethodGet, "/ws", nil)
	assert.True(t, gateway.checkOrigin(reqEmpty))

	reqAllowed, _ := http.NewRequest(http.MethodGet, "/ws", nil)
	reqAllowed.Header.Set("Origin", "http://localhost:5173")
	assert.True(t, gateway.checkOrigin(reqAllowed))

	reqLocalhost, _ := http.NewRequest(http.MethodGet, "/ws", nil)
	reqLocalhost.Header.Set("Origin", "http://127.0.0.1:3000")
	assert.False(t, gateway.checkOrigin(reqLocalhost))

	reqForbidden, _ := http.NewRequest(http.MethodGet, "/ws", nil)
	reqForbidden.Header.Set("Origin", "https://malicious-site.com")
	assert.False(t, gateway.checkOrigin(reqForbidden))

	wildcardGateway := NewPomodoroWSGateway(nil, nil, slog.Default(), "*")
	assert.True(t, wildcardGateway.checkOrigin(reqForbidden))
}

func TestPomodoroWSGateway_HandshakeAndInitialState(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tokenMgr := utils.NewTokenManager("gopa", "test-secret-key-32-chars-long!", time.Hour)
	userID := uuid.New()
	user := domain.User{ID: userID, Email: "ws_test@example.com", DisplayName: "WS Test", Role: domain.RoleUser}

	token, err := tokenMgr.IssueAccessToken(user, uuid.New(), time.Now().UTC())
	require.NoError(t, err)

	now := time.Now().UTC()
	store := &mockPomodoroStore{
		state: &domain.PomodoroState{
			Status:          domain.PomodoroRunning,
			DurationSeconds: 1500,
			StartedAt:       now,
			UpdatedAt:       now,
		},
	}
	pomoService := services.NewPomodoroService(store, mockPomodoroHistoryRepo{})
	gateway := NewPomodoroWSGateway(pomoService, nil, slog.Default(), "*")

	router := gin.New()
	api := router.Group("/api/v1")
	api.GET("/ws/pomodoro", testAuthMiddleware(tokenMgr), func(c *gin.Context) {
		uid := userID
		conn, upgradeErr := gateway.upgrader.Upgrade(c.Writer, c.Request, nil)
		if upgradeErr != nil {
			return
		}
		defer conn.Close()

		st, _ := pomoService.Get(c.Request.Context(), uid)
		_ = conn.WriteJSON(WSMessage{
			Type:      "state",
			Data:      st,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
	})

	server := httptest.NewServer(router)
	defer server.Close()

	// 1. Unauthorized attempt (no token)
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/v1/ws/pomodoro"
	_, resp, dialErr := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.Error(t, dialErr)
	if resp != nil {
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	}

	// 2. Authorized attempt via ?token=
	authedURL := wsURL + "?token=" + token
	conn, _, err := websocket.DefaultDialer.Dial(authedURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	var msg WSMessage
	err = conn.ReadJSON(&msg)
	require.NoError(t, err)
	assert.Equal(t, "state", msg.Type)
	assert.NotNil(t, msg.Data)
}
