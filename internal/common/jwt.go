package common

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

// Sets from main.go during initialization
var TokenAuth *jwtauth.JWTAuth
var AccessTokenTTL time.Duration
var RefreshTokenTTL time.Duration

// InitJWT initializes the JWT authentication middleware with the given secret.
func InitJWT(secret string, accessTokenTTL, refreshTokenTTL time.Duration) {
	TokenAuth = jwtauth.New("HS256", []byte(secret), nil)
	AccessTokenTTL = accessTokenTTL
	RefreshTokenTTL = refreshTokenTTL
}

func GenerateToken(userID uint, duration time.Duration) (string, error) {
	// Access token: short-lived
	_, jwtTokenStr, err := TokenAuth.Encode(map[string]any{
		"user": userID,
		"exp":  time.Now().Add(duration).Unix(),
	})
	if err != nil {
		return "", err
	}

	return jwtTokenStr, nil
}

// Returns user ID if token is valid, error otherwise
func ValidateToken(tokenStr string) (uint, error) {
	token, err := TokenAuth.Decode(tokenStr)
	if err != nil {
		return 0, err
	}
	// Check expiration
	exp := token.Expiration()
	if time.Now().After(exp) {
		return 0, errors.New("token has expired")
	}

	// Get user ID from claims
	userID, ok := token.Get("user")
	if !ok {
		return 0, errors.New("invalid token claim: user")
	}

	userIDUint, err := strconv.ParseUint(fmt.Sprintf("%v", userID), 10, 64)
	if err != nil {
		return 0, errors.New("invalid token claim type: user")
	}

	return uint(userIDUint), nil

}
