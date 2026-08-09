package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gopa/internal/core/domain"
)

type AccessClaims struct {
	Role      domain.Role `json:"role"`
	SessionID string      `json:"sid"`
	jwt.RegisteredClaims
}

type TokenManager struct {
	issuer string
	secret []byte
	ttl    time.Duration
}

type RefreshToken struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
	Secret    string
}

func NewTokenManager(issuer, secret string, ttl time.Duration) *TokenManager {
	return &TokenManager{issuer: issuer, secret: []byte(secret), ttl: ttl}
}

func (m *TokenManager) IssueAccessToken(user domain.User, sessionID uuid.UUID, now time.Time) (string, error) {
	claims := AccessClaims{
		Role:      user.Role,
		SessionID: sessionID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}
	return token, nil
}

func (m *TokenManager) ParseAccessToken(rawToken string) (AccessClaims, error) {
	claims := AccessClaims{}
	token, err := jwt.ParseWithClaims(rawToken, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return m.secret, nil
	}, jwt.WithIssuer(m.issuer))
	if err != nil || !token.Valid {
		return AccessClaims{}, fmt.Errorf("parse access token: %w", err)
	}
	return claims, nil
}

func NewRefreshToken(userID, sessionID uuid.UUID) (RefreshToken, string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return RefreshToken{}, "", fmt.Errorf("generate refresh token: %w", err)
	}
	refreshToken := RefreshToken{UserID: userID, SessionID: sessionID, Secret: base64.RawURLEncoding.EncodeToString(bytes)}
	return refreshToken, refreshToken.String(), nil
}

func ParseRefreshToken(rawToken string) (RefreshToken, error) {
	parts := strings.Split(rawToken, ".")
	if len(parts) != 3 || parts[2] == "" {
		return RefreshToken{}, fmt.Errorf("invalid refresh token")
	}
	userID, err := uuid.Parse(parts[0])
	if err != nil {
		return RefreshToken{}, fmt.Errorf("parse refresh token user ID: %w", err)
	}
	sessionID, err := uuid.Parse(parts[1])
	if err != nil {
		return RefreshToken{}, fmt.Errorf("parse refresh token session ID: %w", err)
	}
	return RefreshToken{UserID: userID, SessionID: sessionID, Secret: parts[2]}, nil
}

func (t RefreshToken) String() string {
	return t.UserID.String() + "." + t.SessionID.String() + "." + t.Secret
}

func HashToken(rawToken string) string {
	hash := sha256.Sum256([]byte(rawToken))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}
