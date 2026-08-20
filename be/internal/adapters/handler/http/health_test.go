package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopa/internal/core/domain"
	"gopa/internal/services"
	"gopa/pkg/utils"
)

type readinessStub struct{ err error }

func (s readinessStub) Check(context.Context) error { return s.err }

func TestHealthHandler_LivenessDoesNotRequireDependencies(t *testing.T) {
	router := newHealthTestRouter(readinessStub{err: errors.New("offline")})
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
}

func TestHealthHandler_ReadinessReturnsServiceUnavailableWhenDependencyFails(t *testing.T) {
	router := newHealthTestRouter(readinessStub{err: errors.New("offline")})
	request := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, response.Code)
	}
}

func newHealthTestRouter(readiness ReadinessChecker) *gin.Engine {
	return NewRouter(
		slog.Default(),
		"http://localhost:5173",
		NewHealthHandler(readiness),
		NewAuthHandler(authServiceStub{}, time.Hour, false),
		nil,
		nil,
		nil,
		nil,
		nil,
		utils.NewTokenManager("gopa", "a secure test secret with enough characters", time.Hour),
		rateLimiterStub{},
	)
}

type authServiceStub struct{}

func (authServiceStub) Register(context.Context, services.RegisterInput) (domain.User, error) {
	return domain.User{}, nil
}
func (authServiceStub) Login(context.Context, services.LoginInput) (services.AuthResult, error) {
	return services.AuthResult{}, nil
}
func (authServiceStub) Refresh(context.Context, string) (services.AuthResult, error) {
	return services.AuthResult{}, nil
}
func (authServiceStub) Logout(context.Context, string) error               { return nil }
func (authServiceStub) Me(context.Context, uuid.UUID) (domain.User, error) { return domain.User{}, nil }
func (authServiceStub) UpdateProfile(context.Context, uuid.UUID, services.UpdateProfileInput) (domain.User, error) {
	return domain.User{}, nil
}

type rateLimiterStub struct{}

func (rateLimiterStub) Allow(context.Context, string, int, time.Duration) (bool, error) {
	return true, nil
}
