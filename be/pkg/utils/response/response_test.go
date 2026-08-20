package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopa/internal/core/domain"
	"gopa/pkg/constants"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(constants.RequestIDKey, "test-req-123")
	return c, w
}

func TestOK(t *testing.T) {
	c, w := setupTestContext()
	data := map[string]string{"message": "success"}

	OK(c, data)

	assert.Equal(t, http.StatusOK, w.Code)
	var env Envelope
	err := json.Unmarshal(w.Body.Bytes(), &env)
	require.NoError(t, err)
	assert.Equal(t, "test-req-123", env.Meta.RequestID)
	assert.Nil(t, env.Error)
	assert.NotNil(t, env.Data)
}

func TestCreated(t *testing.T) {
	c, w := setupTestContext()
	data := map[string]string{"id": "entity-1"}

	Created(c, data)

	assert.Equal(t, http.StatusCreated, w.Code)
	var env Envelope
	err := json.Unmarshal(w.Body.Bytes(), &env)
	require.NoError(t, err)
	assert.Equal(t, "test-req-123", env.Meta.RequestID)
	assert.Nil(t, env.Error)
}

func TestPaginated(t *testing.T) {
	c, w := setupTestContext()
	items := []string{"item1", "item2"}
	total := 100
	cursor := "cur_abc"

	Paginated(c, items, cursor, &total)

	assert.Equal(t, http.StatusOK, w.Code)
	var env Envelope
	err := json.Unmarshal(w.Body.Bytes(), &env)
	require.NoError(t, err)
	assert.Equal(t, "cur_abc", env.Meta.NextCursor)
	assert.Equal(t, 100, *env.Meta.Total)
}

func TestFailDomainErrorMapping(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode int
		expectedAPI  string
	}{
		{"Validation Error", domain.ErrValidation, http.StatusBadRequest, "VALIDATION_ERROR"},
		{"Unauthorized Error", domain.ErrUnauthorized, http.StatusUnauthorized, "UNAUTHORIZED"},
		{"Forbidden Error", domain.ErrForbidden, http.StatusForbidden, "FORBIDDEN"},
		{"Not Found Error", domain.ErrNotFound, http.StatusNotFound, "NOT_FOUND"},
		{"Conflict Error", domain.ErrConflict, http.StatusConflict, "CONFLICT"},
		{"Bad Gateway Error", domain.ErrBadGateway, http.StatusBadGateway, "SERVICE_UNAVAILABLE"},
		{"Unknown Error", errors.New("database explosion"), http.StatusInternalServerError, "INTERNAL_ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := setupTestContext()
			Fail(c, tt.err)

			assert.Equal(t, tt.expectedCode, w.Code)
			var env Envelope
			err := json.Unmarshal(w.Body.Bytes(), &env)
			require.NoError(t, err)
			require.NotNil(t, env.Error)
			assert.Equal(t, tt.expectedAPI, env.Error.Code)
			assert.Equal(t, "test-req-123", env.Meta.RequestID)
		})
	}
}

func TestFailWithDetails(t *testing.T) {
	c, w := setupTestContext()
	details := []ErrDetail{
		{Field: "email", Message: "invalid email format"},
	}

	FailWithDetails(c, "validation failed", details)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var env Envelope
	err := json.Unmarshal(w.Body.Bytes(), &env)
	require.NoError(t, err)
	require.NotNil(t, env.Error)
	assert.Equal(t, "VALIDATION_ERROR", env.Error.Code)
	assert.Len(t, env.Error.Details, 1)
	assert.Equal(t, "email", env.Error.Details[0].Field)
}
