package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopa/internal/core/domain"
	"gopa/internal/core/ports"
	"gopa/pkg/constants"
	"gopa/pkg/utils"
	"gopa/pkg/utils/response"
)

type Identity struct {
	UserID uuid.UUID
}

func Authenticate(tokens *utils.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string
		header := c.GetHeader("Authorization")
		if header != "" {
			parts := strings.SplitN(header, " ", 2)
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				tokenStr = strings.TrimSpace(parts[1])
			}
		}
		if tokenStr == "" && c.Request.URL.Path == "/api/v1/ws/pomodoro" {
			// WebSocket clients cannot send Authorization headers; fall back to query param.
			if t := strings.TrimSpace(c.Query("token")); t != "" {
				tokenStr = t
			} else {
				tokenStr = strings.TrimSpace(c.Query("access_token"))
			}
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
		c.Set(constants.IdentityKey, Identity{UserID: userID})
		c.Next()
	}
}

func rateLimit(limiter ports.RateLimiter, operation string, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		allowed, err := limiter.Allow(c.Request.Context(), operation+":"+c.ClientIP(), limit, window)
		if err != nil {
			response.FailWithStatus(c, http.StatusServiceUnavailable, "RATE_LIMIT_UNAVAILABLE", "Authentication is temporarily unavailable. Please try again later.")
			c.Abort()
			return
		}
		if !allowed {
			response.FailWithStatus(c, http.StatusTooManyRequests, "RATE_LIMITED", "Too many requests. Please try again later.")
			c.Abort()
			return
		}
		c.Next()
	}
}

func identityFromContext(c *gin.Context) (Identity, bool) {
	value, exists := c.Get(constants.IdentityKey)
	if !exists {
		return Identity{}, false
	}
	identity, ok := value.(Identity)
	return identity, ok
}
