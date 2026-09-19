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

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Security-Policy", "default-src 'none'; base-uri 'none'; frame-ancestors 'none'")
		c.Header("Permissions-Policy", "camera=(), geolocation=(), microphone=()")
		c.Header("Referrer-Policy", "no-referrer")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
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
		if origin == "*" {
			// Wildcard: allow all origins (dev only — Validate() blocks this in production).
			// Cannot combine wildcard with Allow-Credentials per the CORS spec.
			c.Header("Access-Control-Allow-Origin", "*")
		} else if requestOrigin != "" && requestOrigin == origin {
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
