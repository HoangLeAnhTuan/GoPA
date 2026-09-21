package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
	const maxJSONBodyBytes = 1 << 20
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxJSONBodyBytes)
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		response.Fail(c, domain.ErrValidation)
		return false
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		response.Fail(c, domain.ErrValidation)
		return false
	}
	if err := binding.Validator.ValidateStruct(target); err != nil {
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

// FlexibleDate provides custom JSON unmarshaling that accepts both HTML date input format
// ("2006-01-02") and standard ISO/RFC3339 timestamps ("2006-01-02T15:04:05Z07:00").
type FlexibleDate time.Time

func (fd *FlexibleDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" || s == "null" {
		return nil
	}
	layouts := []string{
		"2006-01-02",
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02 15:04:05",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			*fd = FlexibleDate(t.UTC())
			return nil
		}
	}
	return errors.New("invalid date format: expected YYYY-MM-DD or RFC3339")
}

func (fd *FlexibleDate) Time() *time.Time {
	if fd == nil {
		return nil
	}
	t := time.Time(*fd)
	return &t
}
