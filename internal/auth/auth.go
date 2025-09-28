package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// JWTClaims represents the claims for both access and refresh tokens
// * It is better to separate them, when refresh token appears.
// * Also Generic Generate and Validate methods will be needed
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

type Auth struct {
	secret []byte
}

// Returns new Auth instance saving secret key
func NewAuth(secret string) *Auth {
	return &Auth{secret: []byte(secret)}
}

// Returns JWT token with given claims
func (a *Auth) GenerateToken(claims *JWTClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, *claims)
	return token.SignedString(a.secret)
}

// ValidateToken validates the JWT token and returns the user ID if valid
func (a *Auth) ValidateToken(tokenString string) (JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return a.secret, nil
	})

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return *claims, nil
	} else {
		return JWTClaims{}, err
	}
}
