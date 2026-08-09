package services

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"gopa/internal/core/domain"
	"gopa/internal/core/ports"
	"gopa/pkg/utils"
)

type RegisterInput struct {
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthResult struct {
	User         domain.User
	AccessToken  string
	RefreshToken string
}

type AuthService struct {
	users           ports.UserRepository
	refreshSessions ports.RefreshSessionStore
	tokens          *utils.TokenManager
	bcryptCost      int
	refreshTokenTTL time.Duration
	now             func() time.Time
}

func NewAuthService(users ports.UserRepository, refreshSessions ports.RefreshSessionStore, tokens *utils.TokenManager, bcryptCost int, refreshTokenTTL time.Duration) *AuthService {
	return &AuthService{
		users:           users,
		refreshSessions: refreshSessions,
		tokens:          tokens,
		bcryptCost:      bcryptCost,
		refreshTokenTTL: refreshTokenTTL,
		now:             time.Now,
	}
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (domain.User, error) {
	email, err := normalizeEmail(input.Email)
	if err != nil {
		return domain.User{}, err
	}
	if err := validatePassword(input.Password); err != nil {
		return domain.User{}, err
	}
	hash, err := utils.HashPassword(input.Password, s.bcryptCost)
	if err != nil {
		return domain.User{}, err
	}
	now := s.now().UTC()
	user := domain.User{ID: uuid.New(), Email: email, PasswordHash: hash, Role: domain.RoleUser, CreatedAt: now, UpdatedAt: now}
	created, err := s.users.Create(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}
	return created, nil
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (AuthResult, error) {
	email, err := normalizeEmail(input.Email)
	if err != nil {
		return AuthResult{}, domain.ErrUnauthorized
	}
	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return AuthResult{}, domain.ErrUnauthorized
		}
		return AuthResult{}, fmt.Errorf("find login user: %w", err)
	}
	if err := utils.VerifyPassword(user.PasswordHash, input.Password); err != nil {
		return AuthResult{}, domain.ErrUnauthorized
	}
	return s.createSession(ctx, user)
}

func (s *AuthService) Refresh(ctx context.Context, rawRefreshToken string) (AuthResult, error) {
	oldToken, err := utils.ParseRefreshToken(rawRefreshToken)
	if err != nil {
		return AuthResult{}, domain.ErrUnauthorized
	}
	user, err := s.users.FindByID(ctx, oldToken.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return AuthResult{}, domain.ErrUnauthorized
		}
		return AuthResult{}, fmt.Errorf("find refresh user: %w", err)
	}
	newSessionID := uuid.New()
	_, newRawToken, err := utils.NewRefreshToken(user.ID, newSessionID)
	if err != nil {
		return AuthResult{}, err
	}
	rotated, err := s.refreshSessions.Rotate(ctx, oldToken.UserID, oldToken.SessionID, utils.HashToken(rawRefreshToken), newSessionID, utils.HashToken(newRawToken), s.refreshTokenTTL)
	if err != nil {
		return AuthResult{}, err
	}
	if !rotated {
		return AuthResult{}, domain.ErrUnauthorized
	}
	accessToken, err := s.tokens.IssueAccessToken(user, newSessionID, s.now().UTC())
	if err != nil {
		return AuthResult{}, err
	}
	return AuthResult{User: user, AccessToken: accessToken, RefreshToken: newRawToken}, nil
}

func (s *AuthService) Logout(ctx context.Context, rawRefreshToken string) error {
	refreshToken, err := utils.ParseRefreshToken(rawRefreshToken)
	if err != nil {
		return nil
	}
	if err := s.refreshSessions.Delete(ctx, refreshToken.UserID, refreshToken.SessionID); err != nil {
		return err
	}
	return nil
}

func (s *AuthService) Me(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return domain.User{}, fmt.Errorf("get current user: %w", err)
	}
	return user, nil
}

func (s *AuthService) createSession(ctx context.Context, user domain.User) (AuthResult, error) {
	sessionID := uuid.New()
	_, rawRefreshToken, err := utils.NewRefreshToken(user.ID, sessionID)
	if err != nil {
		return AuthResult{}, err
	}
	if err := s.refreshSessions.Save(ctx, user.ID, sessionID, utils.HashToken(rawRefreshToken), s.refreshTokenTTL); err != nil {
		return AuthResult{}, err
	}
	accessToken, err := s.tokens.IssueAccessToken(user, sessionID, s.now().UTC())
	if err != nil {
		return AuthResult{}, err
	}
	return AuthResult{User: user, AccessToken: accessToken, RefreshToken: rawRefreshToken}, nil
}

func normalizeEmail(rawEmail string) (string, error) {
	email := strings.ToLower(strings.TrimSpace(rawEmail))
	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email || len(email) > 254 {
		return "", domain.ErrValidation
	}
	return email, nil
}

func validatePassword(password string) error {
	if len(password) < 8 || len(password) > 72 {
		return domain.ErrValidation
	}
	return nil
}
