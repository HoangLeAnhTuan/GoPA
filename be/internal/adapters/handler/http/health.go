package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gopa/pkg/utils/response"
)

type ReadinessChecker interface {
	Check(context.Context) error
}

type HealthHandler struct {
	readiness ReadinessChecker
}

func NewHealthHandler(readiness ReadinessChecker) *HealthHandler {
	return &HealthHandler{readiness: readiness}
}

func (h *HealthHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/healthz", h.liveness)
	router.GET("/readyz", h.ready)
}

func (h *HealthHandler) liveness(c *gin.Context) {
	response.OK(c, gin.H{"status": "ok"})
}

func (h *HealthHandler) ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()
	if err := h.readiness.Check(ctx); err != nil {
		response.FailWithStatus(c, http.StatusServiceUnavailable, "NOT_READY", "Required dependencies are unavailable.")
		return
	}
	response.OK(c, gin.H{"status": "ready"})
}
