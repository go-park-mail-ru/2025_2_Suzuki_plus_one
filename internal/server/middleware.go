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
// Set context [server.AuthContextKey] whether token is valid or not
func authMiddleware(next http.Handler, secret string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authenticator := auth.NewAuth(secret)
		var NoTokenFound bool
		var TokenFailed bool
		var JWTClaims auth.JWTClaims

		// Check token existence and validity
		// Update context fields
		if authenticator.TokenMgr.Exists(r) == false {
			NoTokenFound = true
		} else {
			token := authenticator.TokenMgr.Get(r)
			if token == "" {
				TokenFailed = true
			} else {
				claims, err := authenticator.ValidateToken(token)
				if err != nil {
					log.Println("authMiddleware: error parsing token:", err)
					TokenFailed = true
				}
				JWTClaims = claims
			}
		}

		// Set authContext in request context
		ctx := context.WithValue(r.Context(), AuthContextKey, &authContext{
			Claims:       JWTClaims,
			TokenFailed:  TokenFailed,
			NoTokenFound: NoTokenFound,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Wrapper that ensures authentication is valid
// If not, returns 401 Unauthorized and JSON error [models.ErrorResponse]
func withAuthRequired(nextFunc func(http.ResponseWriter, *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		val := r.Context().Value(AuthContextKey)
		if val == nil {
			log.Println("withAuthRequired: val is nil")
			responseWithError(w, http.StatusUnauthorized, ErrInvalidOrExpired)
			return
		}
		authCtx, ok := val.(*authContext)
		if !ok {
			log.Println("withAuthRequired: no auth context found in request context")
			responseWithError(w, http.StatusUnauthorized, ErrInvalidOrExpired)
			return
		}

		if authCtx == nil || authCtx.NoTokenFound {
			responseWithError(w, http.StatusUnauthorized, ErrNoTokenProvided)
			return
		} else if authCtx.TokenFailed {
			responseWithError(w, http.StatusUnauthorized, ErrInvalidOrExpired)
			return
		}
		http.HandlerFunc(nextFunc).ServeHTTP(w, r)
	}
}
