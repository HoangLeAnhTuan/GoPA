package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopa/internal/core/domain"
	"gopa/pkg/utils/response"
)

func userID(c *gin.Context) uuid.UUID {
	identity, _ := identityFromContext(c)
	return identity.UserID
}

func parseParamUUID(c *gin.Context, param string) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param(param))
	if err != nil {
		response.Fail(c, domain.ErrValidation)
		return uuid.Nil, false
	}
	return id, true
}

func bindJSON(c *gin.Context, target any) bool {
	if err := c.ShouldBindJSON(target); err != nil {
		response.Fail(c, domain.ErrValidation)
		return false
	}
	return true
}

func queryLimit(c *gin.Context, defaultLimit int) int {
	valStr := c.DefaultQuery("limit", strconv.Itoa(defaultLimit))
	val, err := strconv.Atoi(valStr)
	if err != nil || val <= 0 {
		return defaultLimit
	}
	if val > 100 {
		return 100
	}
	return val
}
