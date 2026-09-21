package ws

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gopa/internal/adapters/cache"
	"gopa/internal/core/domain"
	"gopa/internal/services"
	"gopa/pkg/constants"
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
	if g.webOrigin == "*" {
		return true
	}
	if origin == "" {
		// Require an Origin header in non-wildcard mode; empty origin means a non-browser
		// client that bypasses origin checking entirely.
		return false
	}
	return g.webOrigin != "" && strings.EqualFold(origin, g.webOrigin)
}

func (g *PomodoroWSGateway) RegisterRoutes(api *gin.RouterGroup, authenticate gin.HandlerFunc) {
	api.GET("/ws/pomodoro", authenticate, g.HandleConnection)
}

func (g *PomodoroWSGateway) HandleConnection(c *gin.Context) {
	val, exists := c.Get(constants.IdentityKey)
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// Identity is set by Authenticate middleware in the http package.
	// We use a local mirror struct to avoid a circular import.
	type identity struct{ UserID uuid.UUID }
	raw, ok := val.(interface{ GetUserID() uuid.UUID })
	var userID uuid.UUID
	if ok {
		userID = raw.GetUserID()
	} else {
		// Fallback: the value stored is the concrete Identity struct from the http package;
		// extract via JSON round-trip to avoid the cross-package import.
		if b, err := json.Marshal(val); err == nil {
			var tmp struct {
				UserID string `json:"UserID"`
			}
			if json.Unmarshal(b, &tmp) == nil {
				userID, _ = uuid.Parse(tmp.UserID)
			}
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
			// Send a real WebSocket ping so the pong handler resets the read deadline.
			// A JSON "heartbeat" message is not a protocol-level keep-alive.
			_ = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if writeErr := conn.WriteMessage(websocket.PingMessage, nil); writeErr != nil {
				return
			}
		}
	}
}
