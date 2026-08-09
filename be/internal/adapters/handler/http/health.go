package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	c.JSON(http.StatusOK, success(c, gin.H{"status": "ok"}))
}

func (h *HealthHandler) ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()
	if err := h.readiness.Check(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, failure(c, "NOT_READY", "Required dependencies are unavailable."))
		return
	}
	c.JSON(http.StatusOK, success(c, gin.H{"status": "ready"}))
}

func success(c *gin.Context, data any) gin.H {
	return gin.H{"data": data, "meta": gin.H{"request_id": requestID(c)}}
}

func failure(c *gin.Context, code, message string) gin.H {
	return gin.H{"error": gin.H{"code": code, "message": message}, "meta": gin.H{"request_id": requestID(c)}}
}
