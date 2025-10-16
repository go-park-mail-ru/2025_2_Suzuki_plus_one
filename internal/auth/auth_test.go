package auth

import (
	"testing"
	"time"

	"go.uber.org/zap/zaptest"
)

func TestGenerateAndValidateToken(t *testing.T) {
	logger := zaptest.NewLogger(t)
	defer logger.Sync()

	secret := "mysecret"
	userID := "user123"
	tokenType := "access"
	expiration := time.Minute * 5
	appName := "testApp"

	auth := NewAuth(secret, logger)
	claims := NewJWTClaims(userID, tokenType, expiration, appName)

	tokenString, err := auth.GenerateToken(claims)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	validatedClaims, err := auth.ValidateToken(tokenString)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if validatedClaims.Subject != userID {
		t.Errorf("Expected userID %s, got %s", userID, validatedClaims.Subject)
	}
	if validatedClaims.Type != tokenType {
		t.Errorf("Expected token type %s, got %s", tokenType, validatedClaims.Type)
	}
	if validatedClaims.Issuer != appName {
		t.Errorf("Expected issuer %s, got %s", appName, validatedClaims.Issuer)
	}
}

func TestValidateToken_InvalidSignature(t *testing.T) {
	logger := zaptest.NewLogger(t)
	defer logger.Sync()

	secret := "mysecret"
	badSecret := "badsecret"
	userID := "user123"
	tokenType := "access"
	expiration := time.Minute * 5
	appName := "testApp"

	auth := NewAuth(secret, logger)
	badAuth := NewAuth(badSecret, logger)
	claims := NewJWTClaims(userID, tokenType, expiration, appName)

	tokenString, err := auth.GenerateToken(claims)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	_, err = badAuth.ValidateToken(tokenString)
	if err == nil {
		t.Error("Expected error for invalid signature, got nil")
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	logger := zaptest.NewLogger(t)
	defer logger.Sync()

	secret := "mysecret"
	userID := "user123"
	tokenType := "access"
	expiration := -time.Minute // already expired
	appName := "testApp"

	auth := NewAuth(secret, logger)
	claims := NewJWTClaims(userID, tokenType, expiration, appName)

	tokenString, err := auth.GenerateToken(claims)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	_, err = auth.ValidateToken(tokenString)
	if err == nil {
		t.Error("Expected error for expired token, got nil")
	}
}

func TestNewJWTClaimsFields(t *testing.T) {
	userID := "user123"
	tokenType := "refresh"
	expiration := time.Hour
	appName := "testApp"

	claims := NewJWTClaims(userID, tokenType, expiration, appName)
	if claims.Subject != userID {
		t.Errorf("Expected Subject %s, got %s", userID, claims.Subject)
	}
	if claims.Type != tokenType {
		t.Errorf("Expected Type %s, got %s", tokenType, claims.Type)
	}
	if claims.Issuer != appName {
		t.Errorf("Expected Issuer %s, got %s", appName, claims.Issuer)
	}
	if claims.ExpiresAt.Time.Before(time.Now()) {
		t.Error("ExpiresAt should be in the future")
	}
}
