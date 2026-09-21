package response

import (
	"errors"
	"net/http"

	"gopa/internal/core/domain"
	"gopa/pkg/constants"

	"github.com/gin-gonic/gin"
)

// Envelope is the standard top-level JSON response container for all API calls.
type Envelope struct {
	Data  any  `json:"data,omitempty"`
	Error *Err `json:"error,omitempty"`
	Meta  Meta `json:"meta"`
}

// Meta contains operational metadata including tracing request ID and pagination markers.
type Meta struct {
	RequestID  string `json:"request_id"`
	NextCursor string `json:"next_cursor,omitempty"`
	Total      *int   `json:"total,omitempty"`
}

// Err contains machine-readable error codes and human-readable messages.
type Err struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details []ErrDetail `json:"details,omitempty"`
}

// ErrDetail provides contextual field-level validation errors.
type ErrDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// OK writes a 200 OK JSON response enveloped with data and request metadata.
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Envelope{
		Data: data,
		Meta: getMeta(c, "", nil),
	})
}

// Created writes a 201 Created JSON response enveloped with created entity and request metadata.
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Envelope{
		Data: data,
		Meta: getMeta(c, "", nil),
	})
}

// Paginated writes a 200 OK JSON response with list data, optional next cursor, and total count.
func Paginated(c *gin.Context, data any, cursor string, total *int) {
	c.JSON(http.StatusOK, Envelope{
		Data: data,
		Meta: getMeta(c, cursor, total),
	})
}

// Fail maps known domain errors to appropriate HTTP status codes and API error envelopes.
func Fail(c *gin.Context, err error) {
	status, code, message := mapDomainError(err)
	writeError(c, status, code, message, nil)
}

// FailWithStatus writes an explicit operational error while preserving the standard envelope.
func FailWithStatus(c *gin.Context, status int, code, message string) {
	writeError(c, status, code, message, nil)
}

func writeError(c *gin.Context, status int, code, message string, details []ErrDetail) {
	c.JSON(status, Envelope{
		Error: &Err{
			Code:    code,
			Message: message,
			Details: details,
		},
		Meta: getMeta(c, "", nil),
	})
}

// FailWithDetails writes a validation failure response with granular field errors.
func FailWithDetails(c *gin.Context, message string, details []ErrDetail) {
	writeError(c, http.StatusBadRequest, "VALIDATION_ERROR", message, details)
}

func mapDomainError(err error) (int, string, string) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		return http.StatusBadRequest, "VALIDATION_ERROR", "One or more fields are invalid."
	case errors.Is(err, domain.ErrUnauthorized):
		return http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required or credentials are invalid."
	case errors.Is(err, domain.ErrForbidden):
		return http.StatusForbidden, "FORBIDDEN", "You do not have permission to access this resource."
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound, "NOT_FOUND", "The requested resource was not found."
	case errors.Is(err, domain.ErrConflict):
		return http.StatusConflict, "CONFLICT", "A record with this value already exists."
	case errors.Is(err, domain.ErrBadGateway):
		return http.StatusBadGateway, "BAD_GATEWAY", "An upstream dependency failed to process the request."
	default:
		return http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred."
	}
}

func getMeta(c *gin.Context, cursor string, total *int) Meta {
	reqID := "unknown"
	if value, exists := c.Get(constants.RequestIDKey); exists {
		if id, ok := value.(string); ok && id != "" {
			reqID = id
		}
	}
	return Meta{
		RequestID:  reqID,
		NextCursor: cursor,
		Total:      total,
	}
}
