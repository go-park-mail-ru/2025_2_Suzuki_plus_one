// Package auth provides methods for generating and validating JWT tokens
// and handling token storage in HTTP requests (headers or cookies).
package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the claims for both access and refresh tokens
// * It is better to separate them, when refresh token appears.
// * Also Generic Generate and Validate methods will be needed
// > However, we can ignore type bc refresh will be http only cookie
type JWTClaims struct {
	jwt.RegisteredClaims
	Type string `json:"type"` // "access" or "refresh"
}

func NewJWTClaims(userID, tokenType string, expiration time.Duration, appName string) *JWTClaims {
	return &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			Issuer:    appName,
			Subject:   userID,
		},
		Type: tokenType,
	}
}

// Token extraction from request interface
type TokenHandler interface {
	// Request hadndlers
	Exists(request *http.Request) bool
	Get(request *http.Request) string
	Set(request *http.Request, token string)

	// Response handlers
	ResponseWithAuth(writer http.ResponseWriter, token string, response models.SignInResponse)
	ResponseWithDeauth(writer http.ResponseWriter)
}

// Auth structure that holds methods for generating and validating JWT tokens
type Auth struct {
	TokenMgr TokenHandler
	secret   []byte
}

// Returns new Auth instance using Authorization header for token storage
func NewAuthHeader(secret string) *Auth {
	token := NewTokenHeader()

	// Returns new Auth instance saving secret key
	return &Auth{
		secret:   []byte(secret),
		TokenMgr: token,
	}
}

// Returns new Auth instance using Cookie for token storage
func NewAuthCookie(secret string) *Auth {
	token := NewTokenCookie()

	// Returns new Auth instance saving secret key
	return &Auth{
		secret:   []byte(secret),
		TokenMgr: token,
	}
}

// Returns new Auth instance using Authorization header for token storage
func NewAuth(secret string) *Auth {
	// Note: change to NewAuthCookie to use cookies instead of headers
	return NewAuthCookie(secret)
	// return NewAuthHeader(secret)
}

// Returns signed token (JWT) from the given claims
func (a *Auth) GenerateToken(claims *JWTClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, *claims)
	return token.SignedString(a.secret)
}

// Checks JWT token in the string.
//
// Returns:
//
//	auth.JWTClaims - if token is valid, otherwise empty claims
//	error - if token is invalid or expired, otherwise nil
func (a *Auth) ValidateToken(tokenString string) (JWTClaims, error) {
	var claims JWTClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		return a.secret, nil
	})

	if err != nil {
		log.Println("ValidateToken: error parsing token: ", token)
		return JWTClaims{}, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return *claims, nil
	} else {
		return JWTClaims{}, err
	}
}
