package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/auth"
)

// Middleware for handling CORS
func corsMiddleware(next http.Handler, frontendOrigin string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", frontendOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Middleware for logging requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("<<< REQUEST: %s %s : from %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
		log.Printf(">>> REPLY: Completed in %v", time.Since(start))
	})
}

// Middleware to always set Content-Type to application/json
func forceJSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// Middleware that checks for authentication
// It parses and validates JWT token from Authorization header
// Puts it to context as [server.AuthContextKey] no matter valid or not
func authMiddleware(next http.Handler, secret string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		// If no token found, set NoTokenFound and continue
		if token == "" {
			ctx := context.WithValue(r.Context(), AuthContextKey, &authContext{
				NoTokenFound: true,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Validate token
		JWTClaims, err := auth.NewAuth(secret).ValidateToken(token)
		ctx := context.WithValue(r.Context(), AuthContextKey, &authContext{
			Claims:       JWTClaims,
			TokenFailed:  err != nil,
			NoTokenFound: false,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Middleware that ensures authentication is valid
// If not, returns 401 Unauthorized and JSON error [models.ErrorResponse]
func withAuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCtx := r.Context().Value(AuthContextKey).(*authContext)

		if authCtx == nil || authCtx.NoTokenFound {
			responseWithError(w, http.StatusUnauthorized, ErrNoTokenProvided)
			return
		} else if authCtx.TokenFailed {
			responseWithError(w, http.StatusUnauthorized, ErrInvalidOrExpired)
			return
		}
		next.ServeHTTP(w, r)
	})
}
