package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopa/internal/core/domain"
	"gopa/internal/core/ports"
	"gopa/pkg/utils"
)

const identityKey = "identity"

type Identity struct {
	UserID uuid.UUID
}

func Authenticate(tokens *utils.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			writeError(c, domain.ErrUnauthorized)
			c.Abort()
			return
		}
		claims, err := tokens.ParseAccessToken(parts[1])
		if err != nil {
			writeError(c, domain.ErrUnauthorized)
			c.Abort()
			return
		}
		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			writeError(c, domain.ErrUnauthorized)
			c.Abort()
			return
		}
		c.Set(identityKey, Identity{UserID: userID})
		c.Next()
	}
}

func rateLimit(limiter ports.RateLimiter, operation string, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		allowed, err := limiter.Allow(c.Request.Context(), operation+":"+c.ClientIP(), limit, window)
		if err != nil || !allowed {
			c.JSON(http.StatusTooManyRequests, failure(c, "RATE_LIMITED", "Too many requests. Please try again later."))
			c.Abort()
			return
		}
		c.Next()
	}
}

func identityFromContext(c *gin.Context) (Identity, bool) {
	value, exists := c.Get(identityKey)
	identity, ok := value.(Identity)
	return identity, exists && ok
}
