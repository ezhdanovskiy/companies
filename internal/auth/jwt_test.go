package auth

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateJWT(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		username string
		wantErr  bool
	}{
		{
			name:     "Valid token generation",
			email:    "test@example.com",
			username: "testuser",
			wantErr:  false,
		},
		{
			name:     "Empty email",
			email:    "",
			username: "testuser",
			wantErr:  false,
		},
		{
			name:     "Empty username",
			email:    "test@example.com",
			username: "",
			wantErr:  false,
		},
		{
			name:     "Both empty",
			email:    "",
			username: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateJWT(tt.email, tt.username)
			
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			
			require.NoError(t, err)
			assert.NotEmpty(t, token)
			
			// Verify token structure
			parsedToken, err := jwt.ParseWithClaims(token, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtKey), nil
			})
			
			require.NoError(t, err)
			assert.True(t, parsedToken.Valid)
			
			claims, ok := parsedToken.Claims.(*JWTClaim)
			require.True(t, ok)
			assert.Equal(t, tt.email, claims.Email)
			assert.Equal(t, tt.username, claims.Username)
			assert.True(t, claims.ExpiresAt > time.Now().Unix())
		})
	}
}

func TestValidateToken(t *testing.T) {
	// Generate a valid token for testing
	validToken, err := GenerateJWT("test@example.com", "testuser")
	require.NoError(t, err)

	// Generate an expired token
	expiredClaims := &JWTClaim{
		Email:    "expired@example.com",
		Username: "expireduser",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(-1 * time.Hour).Unix(),
		},
	}
	expiredTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredToken, err := expiredTokenObj.SignedString([]byte(jwtKey))
	require.NoError(t, err)

	// Generate a token with wrong signing key
	wrongKeyClaims := &JWTClaim{
		Email:    "wrong@example.com",
		Username: "wronguser",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
	}
	wrongKeyTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, wrongKeyClaims)
	wrongKeyToken, err := wrongKeyTokenObj.SignedString([]byte("wrongkey"))
	require.NoError(t, err)

	tests := []struct {
		name    string
		token   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Valid token",
			token:   validToken,
			wantErr: false,
		},
		{
			name:    "Expired token",
			token:   expiredToken,
			wantErr: true,
			errMsg:  "token is expired",
		},
		{
			name:    "Invalid token format",
			token:   "invalid.token.format",
			wantErr: true,
		},
		{
			name:    "Empty token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "Token with wrong signing key",
			token:   wrongKeyToken,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateToken(tt.token)
			
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSetJWTKey(t *testing.T) {
	// Save original key
	originalKey := jwtKey
	defer func() {
		jwtKey = originalKey
	}()

	// Test setting a new key
	newKey := "newSecretKey123"
	SetJWTKey(newKey)
	
	// Generate token with new key
	token, err := GenerateJWT("test@example.com", "testuser")
	require.NoError(t, err)
	
	// Validate with new key (should work)
	err = ValidateToken(token)
	assert.NoError(t, err)
	
	// Change key and try to validate (should fail)
	SetJWTKey("differentKey")
	err = ValidateToken(token)
	assert.Error(t, err)
}
