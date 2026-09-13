package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword_And_VerifyPassword(t *testing.T) {
	password := "CorrectHorseBatteryStaple123!"

	hash, err := HashPassword(password, bcrypt.MinCost)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)

	// Valid match
	err = VerifyPassword(hash, password)
	assert.NoError(t, err)

	// Mismatched password
	err = VerifyPassword(hash, "WrongPassword456!")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "verify password")
}

func TestHashPassword_InvalidCost(t *testing.T) {
	// Bcrypt returns error when cost > MaxCost (31)
	_, err := HashPassword("test-password", bcrypt.MaxCost+1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "hash password")

	// Cost < MinCost is automatically upgraded to DefaultCost by bcrypt stdlib
	hash, err := HashPassword("test-password", 1)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestVerifyPassword_InvalidHash(t *testing.T) {
	err := VerifyPassword("not-a-valid-bcrypt-hash", "password")
	assert.Error(t, err)

	err = VerifyPassword("", "password")
	assert.Error(t, err)
}
