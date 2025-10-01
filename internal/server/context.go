// This files contains context keys and structures
// Which are used in middleware and then handlers too

package server

import "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/auth"

type contextKey string

const AuthContextKey contextKey = "authContext"

type authContext struct {
	Claims       auth.JWTClaims
	TokenFailed  bool
	NoTokenFound bool
}
