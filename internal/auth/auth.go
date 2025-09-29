package auth

import (
	"log"
	"time"

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

type Auth struct {
	secret []byte
}

// Returns new Auth instance saving secret key
func NewAuth(secret string) *Auth {
	return &Auth{secret: []byte(secret)}
}

// Retrieves token string from "Authorization" header
// Expects header in format "Bearer <token>"
// Returns empty string if no token found
func (a *Auth) GetTokenFromHeader(authHeader string) string {
	const bearerPrefix = "Bearer "
	// Trim "Bearer " prefix
	if len(authHeader) > len(bearerPrefix) && authHeader[:len(bearerPrefix)] == bearerPrefix {
		return authHeader[len(bearerPrefix):]
	}
	return ""
}

// Returns JWT token with given claims
func (a *Auth) GenerateToken(claims *JWTClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, *claims)
	return token.SignedString(a.secret)
}

// ValidateToken validates the JWT token and returns the user ID if valid
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
