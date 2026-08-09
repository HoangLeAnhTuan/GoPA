package http

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"gopa/internal/core/ports"
	"gopa/pkg/utils"
)

func NewRouter(logger *slog.Logger, webOrigin string, health *HealthHandler, auth *AuthHandler, modules *ModuleHandler, tokens *utils.TokenManager, limiter ports.RateLimiter) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), RequestID(), RequestLogger(logger), CORS(webOrigin))
	health.RegisterRoutes(router)
	api := router.Group("/api/v1")
	auth.RegisterRoutes(api, Authenticate(tokens), limiter)
	if modules != nil {
		modules.RegisterRoutes(api, Authenticate(tokens))
	}
	return router
}
