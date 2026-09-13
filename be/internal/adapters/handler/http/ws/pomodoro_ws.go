package ws

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gopa/internal/adapters/cache"
	"gopa/internal/core/domain"
	"gopa/internal/services"
)

type WSMessage struct {
	Type      string `json:"type"`
	Data      any    `json:"data,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

type PomodoroWSGateway struct {
	pomodoro  *services.PomodoroService
	store     *cache.PomodoroStore
	logger    *slog.Logger
	webOrigin string
	upgrader  websocket.Upgrader
}

func NewPomodoroWSGateway(
	pomodoro *services.PomodoroService,
	store *cache.PomodoroStore,
	logger *slog.Logger,
	webOrigin string,
) *PomodoroWSGateway {
	gateway := &PomodoroWSGateway{
		pomodoro:  pomodoro,
		store:     store,
		logger:    logger,
		webOrigin: webOrigin,
	}

	gateway.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return gateway.checkOrigin(r)
		},
	}

	return gateway
}

func (g *PomodoroWSGateway) checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}
	if g.webOrigin == "" || g.webOrigin == "*" {
		return true
	}
	if strings.EqualFold(origin, g.webOrigin) {
		return true
	}
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	originHost := u.Hostname()
	return originHost == "localhost" || originHost == "127.0.0.1"
}

func (g *PomodoroWSGateway) RegisterRoutes(api *gin.RouterGroup, authenticate gin.HandlerFunc) {
	api.GET("/ws/pomodoro", authenticate, g.HandleConnection)
}

func (g *PomodoroWSGateway) HandleConnection(c *gin.Context) {
	var userID uuid.UUID
	if val, exists := c.Get("user_id"); exists {
		if id, ok := val.(uuid.UUID); ok {
			userID = id
		}
	}
	if userID == uuid.Nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	conn, err := g.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		if g.logger != nil {
			g.logger.Warn("websocket upgrade failed", "error", err, "user_id", userID)
		}
		return
	}
	defer func() {
		_ = conn.Close()
	}()

	conn.SetReadLimit(512)
	_ = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		_ = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Send initial snapshot
	ctx := c.Request.Context()
	state, err := g.pomodoro.Get(ctx, userID)
	if err != nil && g.logger != nil {
		g.logger.Debug("get initial pomodoro state", "error", err)
	}

	initialMsg := WSMessage{
		Type:      "state",
		Data:      state,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	_ = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err := conn.WriteJSON(initialMsg); err != nil {
		return
	}

	// Subscribe to Redis Pub/Sub channel
	pubsub := g.store.Subscribe(ctx, userID)
	defer func() {
		_ = pubsub.Close()
	}()

	redisChan := pubsub.Channel()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	clientClosed := make(chan struct{})
	go func() {
		defer close(clientClosed)
		for {
			_, _, readErr := conn.ReadMessage()
			if readErr != nil {
				return
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			_ = conn.WriteControl(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "server shutting down"),
				time.Now().Add(time.Second),
			)
			return

		case <-clientClosed:
			return

		case msg, ok := <-redisChan:
			if !ok {
				return
			}
			var data any
			trimmed := strings.TrimSpace(msg.Payload)
			if trimmed != "" && trimmed != "null" {
				var parsed domain.PomodoroState
				if parseErr := json.Unmarshal([]byte(trimmed), &parsed); parseErr == nil {
					data = parsed
				} else {
					data = trimmed
				}
			}

			_ = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if writeErr := conn.WriteJSON(WSMessage{
				Type:      "state",
				Data:      data,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			}); writeErr != nil {
				return
			}

		case <-ticker.C:
			_ = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if writeErr := conn.WriteJSON(WSMessage{
				Type:      "heartbeat",
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			}); writeErr != nil {
				return
			}
		}
	}
}
