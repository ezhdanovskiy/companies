package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetJWTKey(t *testing.T) {
	originalKey := jwtKey
	defer func() {
		jwtKey = originalKey
	}()

	newKey := "test-secret-key"
	SetJWTKey(newKey)
	assert.Equal(t, newKey, jwtKey)
}

func TestGenerateJWT(t *testing.T) {
	email := "test@example.com"
	username := "testuser"

	token, err := GenerateJWT(email, username)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Parse the token to verify its contents
	parsedToken, err := jwt.ParseWithClaims(token, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	require.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(*JWTClaim)
	require.True(t, ok)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, username, claims.Username)
	assert.Greater(t, claims.ExpiresAt, time.Now().Unix())
}

func TestValidateToken_Valid(t *testing.T) {
	email := "test@example.com"
	username := "testuser"

	token, err := GenerateJWT(email, username)
	require.NoError(t, err)

	err = ValidateToken(token)
	assert.NoError(t, err)
}

func TestValidateToken_Invalid(t *testing.T) {
	err := ValidateToken("invalid.token.here")
	assert.Error(t, err)
}

func TestValidateToken_WrongSignature(t *testing.T) {
	// Create a token with a different key
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTClaim{
		Email:    "test@example.com",
		Username: "testuser",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("wrong-key"))
	require.NoError(t, err)

	err = ValidateToken(tokenString)
	assert.Error(t, err)
}

func TestValidateToken_Expired(t *testing.T) {
	// Create an expired token
	expirationTime := time.Now().Add(-1 * time.Hour) // Already expired
	claims := &JWTClaim{
		Email:    "test@example.com",
		Username: "testuser",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))
	require.NoError(t, err)

	err = ValidateToken(tokenString)
	assert.Error(t, err)
	// Check that error message contains "expired"
	assert.True(t, strings.Contains(err.Error(), "expired"))
}

func TestValidateToken_EmptyFields(t *testing.T) {
	// Create a valid token but with empty Email and Username fields
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTClaim{
		Email:    "",
		Username: "",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))
	require.NoError(t, err)

	// The current implementation doesn't validate if Email/Username are empty,
	// so this should pass validation
	err = ValidateToken(tokenString)
	assert.NoError(t, err)
}
