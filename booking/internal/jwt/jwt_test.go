package jwt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	jwtUtil := NewJWTUtil("test-secret")

	token, err := jwtUtil.GenerateToken("user123", "user", "email@mail.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateToken(t *testing.T) {
	jwtUtil := NewJWTUtil("test-secret")

	token, err := jwtUtil.GenerateToken("user123", "admin", "email@mail.com")
	assert.NoError(t, err)

	claims, err := jwtUtil.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "user123", claims.UserID)
	assert.Equal(t, "admin", claims.Role)
}

func TestValidateInvalidToken(t *testing.T) {
	jwtUtil := NewJWTUtil("test-secret")

	_, err := jwtUtil.ValidateToken("invalid-token")
	assert.Error(t, err)
}
