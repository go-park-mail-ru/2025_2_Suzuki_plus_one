package middleware

import (
	"net/http"
)

// Middleware for handling CORS
func SetCors(next http.Handler, frontendOrigin string) http.Handler {
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
