package http

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopa/internal/core/domain"
	"gopa/internal/core/ports"
	"gopa/internal/services"
)

const refreshCookieName = "gopa_refresh"

type AuthApplication interface {
	Register(context.Context, services.RegisterInput) (domain.User, error)
	Login(context.Context, services.LoginInput) (services.AuthResult, error)
	Refresh(context.Context, string) (services.AuthResult, error)
	Logout(context.Context, string) error
	Me(context.Context, uuid.UUID) (domain.User, error)
}

type AuthHandler struct {
	service       AuthApplication
	refreshTTL    time.Duration
	secureCookies bool
}

func NewAuthHandler(service AuthApplication, refreshTTL time.Duration, secureCookies bool) *AuthHandler {
	return &AuthHandler{service: service, refreshTTL: refreshTTL, secureCookies: secureCookies}
}

func (h *AuthHandler) RegisterRoutes(api *gin.RouterGroup, authenticate gin.HandlerFunc, limiter ports.RateLimiter) {
	auth := api.Group("/auth")
	auth.POST("/register", h.register)
	auth.POST("/login", rateLimit(limiter, "login", 5, time.Minute), h.login)
	auth.POST("/refresh", rateLimit(limiter, "refresh", 10, time.Minute), h.refresh)
	auth.POST("/logout", h.logout)
	api.GET("/me", authenticate, h.me)
}

func (h *AuthHandler) register(c *gin.Context) {
	var request registerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		writeError(c, domain.ErrValidation)
		return
	}
	if request.Password != request.ConfirmPassword {
		writeError(c, domain.ErrValidation)
		return
	}
	user, err := h.service.Register(c.Request.Context(), services.RegisterInput{Email: request.Email, Password: request.Password})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, success(c, userResponse(user)))
}

func (h *AuthHandler) login(c *gin.Context) {
	var request loginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		writeError(c, domain.ErrUnauthorized)
		return
	}
	result, err := h.service.Login(c.Request.Context(), services.LoginInput{Email: request.Email, Password: request.Password})
	if err != nil {
		writeError(c, err)
		return
	}
	h.setRefreshCookie(c, result.RefreshToken)
	c.JSON(http.StatusOK, success(c, authResponse(result)))
}

func (h *AuthHandler) refresh(c *gin.Context) {
	rawRefreshToken, err := c.Cookie(refreshCookieName)
	if err != nil {
		writeError(c, domain.ErrUnauthorized)
		return
	}
	result, err := h.service.Refresh(c.Request.Context(), rawRefreshToken)
	if err != nil {
		h.clearRefreshCookie(c)
		writeError(c, err)
		return
	}
	h.setRefreshCookie(c, result.RefreshToken)
	c.JSON(http.StatusOK, success(c, authResponse(result)))
}

func (h *AuthHandler) logout(c *gin.Context) {
	if rawRefreshToken, err := c.Cookie(refreshCookieName); err == nil {
		if err := h.service.Logout(c.Request.Context(), rawRefreshToken); err != nil {
			writeError(c, err)
			return
		}
	}
	h.clearRefreshCookie(c)
	c.JSON(http.StatusOK, success(c, gin.H{}))
}

func (h *AuthHandler) me(c *gin.Context) {
	identity, ok := identityFromContext(c)
	if !ok {
		writeError(c, domain.ErrUnauthorized)
		return
	}
	user, err := h.service.Me(c.Request.Context(), identity.UserID)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, success(c, userResponse(user)))
}

func (h *AuthHandler) setRefreshCookie(c *gin.Context, token string) {
	http.SetCookie(c.Writer, &http.Cookie{Name: refreshCookieName, Value: token, Path: "/api/v1/auth", MaxAge: int(h.refreshTTL.Seconds()), HttpOnly: true, Secure: h.secureCookies, SameSite: http.SameSiteLaxMode})
}

func (h *AuthHandler) clearRefreshCookie(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{Name: refreshCookieName, Value: "", Path: "/api/v1/auth", MaxAge: -1, HttpOnly: true, Secure: h.secureCookies, SameSite: http.SameSiteLaxMode})
}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		c.JSON(http.StatusBadRequest, failure(c, "VALIDATION_ERROR", "One or more fields are invalid."))
	case errors.Is(err, domain.ErrConflict):
		c.JSON(http.StatusConflict, failure(c, "CONFLICT", "A record with this value already exists."))
	case errors.Is(err, domain.ErrNotFound):
		c.JSON(http.StatusNotFound, failure(c, "NOT_FOUND", "The requested resource was not found."))
	case errors.Is(err, domain.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, failure(c, "UNAUTHORIZED", "Authentication is required or credentials are invalid."))
	default:
		c.JSON(http.StatusInternalServerError, failure(c, "INTERNAL_ERROR", "An unexpected error occurred."))
	}
}

type registerRequest struct {
	Email           string `json:"email" binding:"required,email,max=254"`
	Password        string `json:"password" binding:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8,max=72"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email,max=254"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type userDTO struct {
	ID        uuid.UUID   `json:"id"`
	Email     string      `json:"email"`
	Role      domain.Role `json:"role"`
	CreatedAt time.Time   `json:"created_at"`
}

type authDTO struct {
	User        userDTO `json:"user"`
	AccessToken string  `json:"access_token"`
}

func userResponse(user domain.User) userDTO {
	return userDTO{ID: user.ID, Email: user.Email, Role: user.Role, CreatedAt: user.CreatedAt}
}

func authResponse(result services.AuthResult) authDTO {
	return authDTO{User: userResponse(result.User), AccessToken: result.AccessToken}
}
