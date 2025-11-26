package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// Middleware for logging requests
func GetLogging(l logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// We can pass logger into context here actually
			// https://www.kaznacheev.me/posts/en/where-to-place-logger-in-golang/
			ctx := common.GetContext(r)
			log := logger.LoggerWithKey(l, ctx, common.ContextKeyRequestID)

			startBanner := fmt.Sprintf("---------- [%s] START --------------------", r.URL.Path)
			endBanner := fmt.Sprintf("---------- [%s]  END --------------------", r.URL.Path)
			log.Info(startBanner)
			log.Info("Request started",
				log.ToString("method", r.Method),
				log.ToString("path", r.URL.Path),
				log.ToString("remote_addr", r.RemoteAddr),
				log.ToString("user_agent", r.UserAgent()),
			)

			next.ServeHTTP(w, r)
			duration := time.Since(start)

			log.Info("Request completed",
				log.ToString("method", r.Method),
				log.ToString("path", r.URL.Path),
				log.ToDuration("duration", duration),
				log.ToString("remote_addr", r.RemoteAddr))
			log.Info(endBanner)
		})
	}
}
