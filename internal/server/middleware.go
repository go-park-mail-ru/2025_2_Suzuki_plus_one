package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/auth"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/utils"
	"go.uber.org/zap"
)

// Middleware for handling CORS
func corsMiddleware(next http.Handler, frontendOrigin string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", frontendOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// Allow cookies to be sent in cross-origin requests
		// w.Header().Set("Access-Control-Allow-Credentials", "true") // Requires Allow-origin not '*'
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Middleware for logging requests
func loggingMiddleware(next http.Handler, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		logger.Info("Request started",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("user_agent", r.UserAgent()),
		)

		next.ServeHTTP(w, r)
		duration := time.Since(start)

		logger.Info("Request completed",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Duration("duration", duration),
			zap.String("remote_addr", r.RemoteAddr))
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
func authMiddleware(next http.Handler, secret string, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Auth Middleware",
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method))
		// Get token from Authorization header
		authenticator := auth.NewAuth(secret, logger)
		var NoTokenFound bool
		var TokenFailed bool
		var JWTClaims auth.JWTClaims

		// Check token existence and validity
		// Update context fields
		if authenticator.TokenMgr.Exists(r) == false {
			NoTokenFound = true
			logger.Debug("No token found")
		} else {
			token := authenticator.TokenMgr.Get(r)
			if token == "" {
				TokenFailed = true
				logger.Debug("Empty token found")
			} else {
				claims, err := authenticator.ValidateToken(token)
				if err != nil {
					logger.Warn("Token validation failed",
						zap.String("token_prefix", utils.SafeTokenPrefix(token)),
						zap.Error(err))
					log.Println("authMiddleware: error parsing token:", err)

					TokenFailed = true
				}
				JWTClaims = claims

				logger.Debug("Token validation successful",
					zap.String("subject", claims.Subject),
					zap.String("type", claims.Type))
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
func withAuthRequired(nextFunc func(http.ResponseWriter, *http.Request), logger *zap.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Auth required check",
			zap.String("path", r.URL.Path))

		val := r.Context().Value(AuthContextKey)
		if val == nil {
			logger.Warn("withAuthRequired: val is nil")

			responseWithError(w, http.StatusUnauthorized, ErrInvalidOrExpired, logger)
			return
		}
		authCtx, ok := val.(*authContext)
		if !ok {
			logger.Error("Invalid context type")

			responseWithError(w, http.StatusUnauthorized, ErrInvalidOrExpired, logger)
			return
		}

		if authCtx == nil || authCtx.NoTokenFound {
			logger.Warn("No token provided for protected route")

			responseWithError(w, http.StatusUnauthorized, ErrNoTokenProvided, logger)
			return
		} else if authCtx.TokenFailed {
			logger.Warn("Invalid or expired token")

			responseWithError(w, http.StatusUnauthorized, ErrInvalidOrExpired, logger)
			return
		}
		http.HandlerFunc(nextFunc).ServeHTTP(w, r)
	}
}
