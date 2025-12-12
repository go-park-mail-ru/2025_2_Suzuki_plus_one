package middleware

import (
	"net/http"
	"strings"
)

// Middleware for handling CORS, matching smiddleware signature
func SetCors(frontendOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var frontendOrigins []string
			if frontendOrigin == "" {
				// If no frontend origin is specified, allow all origins
				frontendOrigins = []string{"*"}
			}

			// Split the supplied frontendOrigin string on commas and trim whitespace
			origins := strings.Split(frontendOrigin, ",")
			for _, o := range origins {
				if s := strings.TrimSpace(o); s != "" {
					frontendOrigins = append(frontendOrigins, s)
				}
			}
			// Set CORS headers based on allowed origins
			originHeader := r.Header.Get("Origin")
			allowed := false
			for _, origin := range frontendOrigins {
				if origin == "*" {
					w.Header().Set("Access-Control-Allow-Origin", "*")
					allowed = true
					break
				}
				if origin == originHeader {
					w.Header().Set("Access-Control-Allow-Origin", originHeader)
					allowed = true
					break
				}
			}
			if allowed {
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Vary", "Origin")
				// Allow cookies to be sent in cross-origin requests
				// w.Header().Set("Access-Control-Allow-Credentials", "true") // Requires Allow-origin not '*'
			}
			w.Header().Set("Access-Control-Allow-Credentials", "true") // Requires Allow-origin not '*'
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
