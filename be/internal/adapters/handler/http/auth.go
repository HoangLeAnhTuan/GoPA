package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopa/internal/core/domain"
	"gopa/internal/core/ports"
	"gopa/internal/services"
	"gopa/pkg/constants"
	"gopa/pkg/utils/response"
)

type AuthApplication interface {
	Register(context.Context, services.RegisterInput) (domain.User, error)
	Login(context.Context, services.LoginInput) (services.AuthResult, error)
	Refresh(context.Context, string) (services.AuthResult, error)
	Logout(context.Context, string) error
	Me(context.Context, uuid.UUID) (domain.User, error)
	UpdateProfile(context.Context, uuid.UUID, services.UpdateProfileInput) (domain.User, error)
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
	api.PATCH("/me", authenticate, h.updateProfile)
}

func (h *AuthHandler) register(c *gin.Context) {
	var request registerRequest
	if !bindJSON(c, &request) {
		return
	}
	if request.Password != request.ConfirmPassword {
		response.Fail(c, domain.ErrValidation)
		return
	}
	user, err := h.service.Register(c.Request.Context(), services.RegisterInput{
		Email:       request.Email,
		Password:    request.Password,
		DisplayName: request.DisplayName,
	})
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, userResponse(user))
}

func (h *AuthHandler) login(c *gin.Context) {
	var request loginRequest
	if !bindJSON(c, &request) {
		return
	}
	result, err := h.service.Login(c.Request.Context(), services.LoginInput{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		response.Fail(c, err)
		return
	}
	h.setRefreshCookie(c, result.RefreshToken)
	response.OK(c, authResponse(result))
}

func (h *AuthHandler) refresh(c *gin.Context) {
	rawRefreshToken, err := c.Cookie(constants.RefreshCookieName)
	if err != nil {
		response.Fail(c, domain.ErrUnauthorized)
		return
	}
	result, err := h.service.Refresh(c.Request.Context(), rawRefreshToken)
	if err != nil {
		h.clearRefreshCookie(c)
		response.Fail(c, err)
		return
	}
	h.setRefreshCookie(c, result.RefreshToken)
	response.OK(c, authResponse(result))
}

func (h *AuthHandler) logout(c *gin.Context) {
	if rawRefreshToken, err := c.Cookie(constants.RefreshCookieName); err == nil {
		if err := h.service.Logout(c.Request.Context(), rawRefreshToken); err != nil {
			response.Fail(c, err)
			return
		}
	}
	h.clearRefreshCookie(c)
	response.OK(c, gin.H{"logged_out": true})
}

func (h *AuthHandler) me(c *gin.Context) {
	identity, ok := identityFromContext(c)
	if !ok {
		response.Fail(c, domain.ErrUnauthorized)
		return
	}
	user, err := h.service.Me(c.Request.Context(), identity.UserID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, userResponse(user))
}

type updateProfileRequest struct {
	DisplayName *string `json:"display_name"`
	AvatarURL   *string `json:"avatar_url"`
}

func (h *AuthHandler) updateProfile(c *gin.Context) {
	identity, ok := identityFromContext(c)
	if !ok {
		response.Fail(c, domain.ErrUnauthorized)
		return
	}
	var req updateProfileRequest
	if !bindJSON(c, &req) {
		return
	}
	user, err := h.service.UpdateProfile(c.Request.Context(), identity.UserID, services.UpdateProfileInput{
		DisplayName: req.DisplayName,
		AvatarURL:   req.AvatarURL,
	})
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, userResponse(user))
}

func (h *AuthHandler) setRefreshCookie(c *gin.Context, token string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     constants.RefreshCookieName,
		Value:    token,
		Path:     "/api/v1/auth",
		MaxAge:   int(h.refreshTTL.Seconds()),
		HttpOnly: true,
		Secure:   h.secureCookies,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *AuthHandler) clearRefreshCookie(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     constants.RefreshCookieName,
		Value:    "",
		Path:     "/api/v1/auth",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.secureCookies,
		SameSite: http.SameSiteLaxMode,
	})
}

type registerRequest struct {
	DisplayName     string `json:"display_name"`
	Email           string `json:"email" binding:"required,email,max=254"`
	Password        string `json:"password" binding:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8,max=72"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email,max=254"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type userDTO struct {
	ID          uuid.UUID   `json:"id"`
	Email       string      `json:"email"`
	DisplayName string      `json:"display_name"`
	AvatarURL   *string     `json:"avatar_url,omitempty"`
	Role        domain.Role `json:"role"`
	CreatedAt   time.Time   `json:"created_at"`
}

type authDTO struct {
	User        userDTO `json:"user"`
	AccessToken string  `json:"access_token"`
}

func userResponse(user domain.User) userDTO {
	return userDTO{
		ID:          user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
	}
}

func authResponse(result services.AuthResult) authDTO {
	return authDTO{User: userResponse(result.User), AccessToken: result.AccessToken}
}
