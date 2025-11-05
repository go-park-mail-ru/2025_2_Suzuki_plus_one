package middleware

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// Middleware for logging requests
func GetLogging(logger logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// We can pass logger into context here if needed
			// https://www.kaznacheev.me/posts/en/where-to-place-logger-in-golang/

			logger.Info("Request started",
				logger.ToString("method", r.Method),
				logger.ToString("path", r.URL.Path),
				logger.ToString("remote_addr", r.RemoteAddr),
				logger.ToString("user_agent", r.UserAgent()),
			)

			next.ServeHTTP(w, r)
			duration := time.Since(start)

			logger.Info("Request completed",
				logger.ToString("method", r.Method),
				logger.ToString("path", r.URL.Path),
				logger.ToDuration("duration", duration),
				logger.ToString("remote_addr", r.RemoteAddr))
		})
	}
}
