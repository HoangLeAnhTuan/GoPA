package http

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"gopa/internal/adapters/handler/http/ws"
	"gopa/internal/core/ports"
	"gopa/pkg/utils"
)

func NewRouter(
	logger *slog.Logger,
	webOrigin string,
	health *HealthHandler,
	auth *AuthHandler,
	tasks *TasksHandler,
	linguistics *LinguisticsHandler,
	finance *FinanceHandler,
	journal *JournalHandler,
	pomodoro *PomodoroHandler,
	pomodoroWS *ws.PomodoroWSGateway,
	tokens *utils.TokenManager,
	limiter ports.RateLimiter,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), RequestID(), SecurityHeaders(), RequestLogger(logger), CORS(webOrigin))
	if health != nil {
		health.RegisterRoutes(router)
	}
	api := router.Group("/api/v1")
	authMiddleware := Authenticate(tokens)
	if auth != nil {
		auth.RegisterRoutes(api, authMiddleware, limiter)
	}
	if tasks != nil {
		tasks.RegisterRoutes(api, authMiddleware)
	}
	if linguistics != nil {
		linguistics.RegisterRoutes(api, authMiddleware)
	}
	if finance != nil {
		finance.RegisterRoutes(api, authMiddleware)
	}
	if journal != nil {
		journal.RegisterRoutes(api, authMiddleware)
	}
	if pomodoro != nil {
		pomodoro.RegisterRoutes(api, authMiddleware)
	}
	if pomodoroWS != nil {
		pomodoroWS.RegisterRoutes(api, authMiddleware)
	}
	return router
}
