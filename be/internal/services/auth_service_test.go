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
func (s *userRepositoryStub) Update(_ context.Context, user domain.User) (domain.User, error) {
	s.user = user
	return user, nil
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

	user, err := service.Register(context.Background(), RegisterInput{
		Email:       " User@Example.COM ",
		Password:    "password-123",
		DisplayName: "  Test User  ",
	})

	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if user.Email != "user@example.com" {
		t.Fatalf("expected canonical email, got %q", user.Email)
	}
	if user.DisplayName != "Test User" {
		t.Fatalf("expected trimmed display name, got %q", user.DisplayName)
	}
	if err := utils.VerifyPassword(user.PasswordHash, "password-123"); err != nil {
		t.Fatalf("password was not hashed correctly: %v", err)
	}
}

func TestAuthService_RegisterDisplayNameTooLong(t *testing.T) {
	service := NewAuthService(&userRepositoryStub{}, &refreshSessionStub{}, utils.NewTokenManager("gopa", "a secure development secret with enough characters", time.Minute), bcrypt.MinCost, time.Hour)

	longName := string(make([]byte, 101))
	_, err := service.Register(context.Background(), RegisterInput{
		Email:       "user@example.com",
		Password:    "password-123",
		DisplayName: longName,
	})

	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected validation error for long display name, got %v", err)
	}
}

func TestAuthService_UpdateProfile(t *testing.T) {
	userID := uuid.New()
	existingUser := domain.User{
		ID:          userID,
		Email:       "test@example.com",
		DisplayName: "Old Name",
		Role:        domain.RoleUser,
	}
	users := &userRepositoryStub{user: existingUser}
	service := NewAuthService(users, &refreshSessionStub{}, utils.NewTokenManager("gopa", "a secure development secret with enough characters", time.Minute), bcrypt.MinCost, time.Hour)

	// 1. Update DisplayName and AvatarURL
	newName := "New Name"
	avatar := "https://example.com/avatar.png"
	updated, err := service.UpdateProfile(context.Background(), userID, UpdateProfileInput{
		DisplayName: &newName,
		AvatarURL:   &avatar,
	})
	if err != nil {
		t.Fatalf("update profile: %v", err)
	}
	if updated.DisplayName != "New Name" {
		t.Fatalf("expected DisplayName 'New Name', got %q", updated.DisplayName)
	}
	if updated.AvatarURL == nil || *updated.AvatarURL != avatar {
		t.Fatalf("expected AvatarURL %q, got %v", avatar, updated.AvatarURL)
	}

	// 2. Clear AvatarURL with empty string
	emptyAvatar := ""
	cleared, err := service.UpdateProfile(context.Background(), userID, UpdateProfileInput{
		AvatarURL: &emptyAvatar,
	})
	if err != nil {
		t.Fatalf("update profile clear avatar: %v", err)
	}
	if cleared.AvatarURL != nil {
		t.Fatalf("expected AvatarURL to be nil, got %v", cleared.AvatarURL)
	}

	// 3. User Not Found
	_, err = service.UpdateProfile(context.Background(), uuid.New(), UpdateProfileInput{
		DisplayName: &newName,
	})
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for nonexistent user, got %v", err)
	}

	// 4. Validation Error on too long DisplayName
	longName := string(make([]byte, 101))
	_, err = service.UpdateProfile(context.Background(), userID, UpdateProfileInput{
		DisplayName: &longName,
	})
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected ErrValidation on long DisplayName, got %v", err)
	}

	// 5. Validation Error on too long AvatarURL
	longAvatar := "https://example.com/" + string(make([]byte, 2040))
	_, err = service.UpdateProfile(context.Background(), userID, UpdateProfileInput{
		AvatarURL: &longAvatar,
	})
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected ErrValidation on long AvatarURL, got %v", err)
	}
}

func TestAuthService_LoginReturnsUnauthorizedForUnknownEmail(t *testing.T) {
	service := NewAuthService(&userRepositoryStub{}, &refreshSessionStub{}, utils.NewTokenManager("gopa", "a secure development secret with enough characters", time.Minute), bcrypt.MinCost, time.Hour)

	_, err := service.Login(context.Background(), LoginInput{Email: "unknown@example.com", Password: "password-123"})

	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("expected unauthorized error, got %v", err)
	}
}
