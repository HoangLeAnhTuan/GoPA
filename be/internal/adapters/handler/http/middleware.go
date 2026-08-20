package http

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"gopa/pkg/constants"

	"github.com/gin-gonic/gin"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := strings.TrimSpace(c.GetHeader(constants.RequestIDHeader))
		if len(requestID) == 0 || len(requestID) > 128 {
			requestID = newRequestID()
		}
		c.Set(constants.RequestIDKey, requestID)
		c.Header(constants.RequestIDHeader, requestID)
		c.Next()
	}
}

func RequestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startedAt := time.Now()
		c.Next()
		logger.Info("http request completed",
			"request_id", requestID(c),
			"method", c.Request.Method,
			"route", c.FullPath(),
			"status", c.Writer.Status(),
			"latency_ms", time.Since(startedAt).Milliseconds(),
		)
	}
}

func CORS(origin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestOrigin := c.GetHeader("Origin")
		if requestOrigin != "" && requestOrigin == origin {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Vary", "Origin")
		}
		if c.Request.Method == http.MethodOptions {
			c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, "+constants.RequestIDHeader)
			c.Status(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func requestID(c *gin.Context) string {
	value, exists := c.Get(constants.RequestIDKey)
	if requestID, ok := value.(string); exists && ok {
		return requestID
	}
	return "unknown"
}

func newRequestID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "unavailable"
	}
	return hex.EncodeToString(bytes)
}
