package utils

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopa/internal/core/domain"
)

func TestTokenManager_IssueAndParseAccessToken(t *testing.T) {
	manager := NewTokenManager("gopa-test", "super-secret-key-123456789012345", 15*time.Minute)
	userID := uuid.New()
	sessionID := uuid.New()
	user := domain.User{
		ID:   userID,
		Role: domain.RoleUser,
	}
	now := time.Now().UTC().Truncate(time.Second)

	token, err := manager.IssueAccessToken(user, sessionID, now)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := manager.ParseAccessToken(token)
	require.NoError(t, err)
	assert.Equal(t, userID.String(), claims.Subject)
	assert.Equal(t, domain.RoleUser, claims.Role)
	assert.Equal(t, sessionID.String(), claims.SessionID)
	assert.Equal(t, "gopa-test", claims.Issuer)
	assert.Equal(t, now.Unix(), claims.IssuedAt.Unix())
	assert.Equal(t, now.Add(15*time.Minute).Unix(), claims.ExpiresAt.Unix())
}

func TestTokenManager_ExpiredToken(t *testing.T) {
	manager := NewTokenManager("gopa-test", "super-secret-key-123456789012345", 100*time.Millisecond)
	userID := uuid.New()
	sessionID := uuid.New()
	user := domain.User{ID: userID, Role: domain.RoleUser}

	// Issue token back in the past
	pastTime := time.Now().UTC().Add(-1 * time.Hour)
	token, err := manager.IssueAccessToken(user, sessionID, pastTime)
	require.NoError(t, err)

	_, err = manager.ParseAccessToken(token)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parse access token")
}

func TestTokenManager_WrongSecret(t *testing.T) {
	manager1 := NewTokenManager("gopa-test", "secret-key-alpha-12345678901234", 15*time.Minute)
	manager2 := NewTokenManager("gopa-test", "secret-key-beta-987654321098765", 15*time.Minute)

	user := domain.User{ID: uuid.New(), Role: domain.RoleAdmin}
	token, err := manager1.IssueAccessToken(user, uuid.New(), time.Now().UTC())
	require.NoError(t, err)

	_, err = manager2.ParseAccessToken(token)
	assert.Error(t, err)
}

func TestTokenManager_WrongIssuer(t *testing.T) {
	manager1 := NewTokenManager("issuer-a", "common-secret-key-1234567890123", 15*time.Minute)
	manager2 := NewTokenManager("issuer-b", "common-secret-key-1234567890123", 15*time.Minute)

	user := domain.User{ID: uuid.New(), Role: domain.RoleUser}
	token, err := manager1.IssueAccessToken(user, uuid.New(), time.Now().UTC())
	require.NoError(t, err)

	_, err = manager2.ParseAccessToken(token)
	assert.Error(t, err)
}

func TestRefreshToken_Lifecycle(t *testing.T) {
	userID := uuid.New()
	sessionID := uuid.New()

	tokenObj, rawToken, err := NewRefreshToken(userID, sessionID)
	require.NoError(t, err)
	assert.Equal(t, userID, tokenObj.UserID)
	assert.Equal(t, sessionID, tokenObj.SessionID)
	assert.NotEmpty(t, tokenObj.Secret)
	assert.Equal(t, tokenObj.String(), rawToken)

	parsed, err := ParseRefreshToken(rawToken)
	require.NoError(t, err)
	assert.Equal(t, userID, parsed.UserID)
	assert.Equal(t, sessionID, parsed.SessionID)
	assert.Equal(t, tokenObj.Secret, parsed.Secret)
}

func TestParseRefreshToken_InvalidFormats(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Empty string", ""},
		{"Single part", "justatoken"},
		{"Two parts", "part1.part2"},
		{"Invalid user UUID", "not-a-uuid." + uuid.New().String() + ".secret"},
		{"Invalid session UUID", uuid.New().String() + ".not-a-uuid.secret"},
		{"Empty secret part", uuid.New().String() + "." + uuid.New().String() + "."},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseRefreshToken(tc.input)
			assert.Error(t, err)
		})
	}
}

func TestHashToken(t *testing.T) {
	token1 := "user1.session1.secret1"
	token2 := "user2.session2.secret2"

	hash1 := HashToken(token1)
	hash1Repeat := HashToken(token1)
	hash2 := HashToken(token2)

	assert.NotEmpty(t, hash1)
	assert.Equal(t, hash1, hash1Repeat, "Hash must be deterministic")
	assert.NotEqual(t, hash1, hash2, "Distinct tokens must have distinct hashes")
}
