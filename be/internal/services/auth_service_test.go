package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gopa/internal/core/domain"
	"gopa/pkg/utils"
)

type userRepositoryStub struct {
	user      domain.User
	createErr error
}

func (s *userRepositoryStub) Create(_ context.Context, user domain.User) (domain.User, error) {
	s.user = user
	return user, s.createErr
}
func (s *userRepositoryStub) FindByID(_ context.Context, id uuid.UUID) (domain.User, error) {
	if s.user.ID == id {
		return s.user, nil
	}
	return domain.User{}, domain.ErrNotFound
}
func (s *userRepositoryStub) FindByEmail(_ context.Context, email string) (domain.User, error) {
	if s.user.Email == email {
		return s.user, nil
	}
	return domain.User{}, domain.ErrNotFound
}

type refreshSessionStub struct {
	saved  bool
	rotate bool
}

func (s *refreshSessionStub) Save(context.Context, uuid.UUID, uuid.UUID, string, time.Duration) error {
	s.saved = true
	return nil
}
func (s *refreshSessionStub) Rotate(context.Context, uuid.UUID, uuid.UUID, string, uuid.UUID, string, time.Duration) (bool, error) {
	return s.rotate, nil
}
func (s *refreshSessionStub) Delete(context.Context, uuid.UUID, uuid.UUID) error { return nil }

func TestAuthService_RegisterCanonicalizesEmailAndHashesPassword(t *testing.T) {
	users := &userRepositoryStub{}
	service := NewAuthService(users, &refreshSessionStub{}, utils.NewTokenManager("gopa", "a secure development secret with enough characters", time.Minute), bcrypt.MinCost, time.Hour)

	user, err := service.Register(context.Background(), RegisterInput{Email: " User@Example.COM ", Password: "password-123"})

	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if user.Email != "user@example.com" {
		t.Fatalf("expected canonical email, got %q", user.Email)
	}
	if err := utils.VerifyPassword(user.PasswordHash, "password-123"); err != nil {
		t.Fatalf("password was not hashed correctly: %v", err)
	}
}

func TestAuthService_LoginReturnsUnauthorizedForUnknownEmail(t *testing.T) {
	service := NewAuthService(&userRepositoryStub{}, &refreshSessionStub{}, utils.NewTokenManager("gopa", "a secure development secret with enough characters", time.Minute), bcrypt.MinCost, time.Hour)

	_, err := service.Login(context.Background(), LoginInput{Email: "unknown@example.com", Password: "password-123"})

	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("expected unauthorized error, got %v", err)
	}
}
