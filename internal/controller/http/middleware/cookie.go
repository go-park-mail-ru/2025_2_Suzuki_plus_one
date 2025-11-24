package middleware

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http/handlers"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

var (
	ErrorMissingCookie = handlers.ResponseError{
		Code:    http.StatusUnauthorized,
		Message: errors.New("Missing required cookie"),
	}
)

// MustCookieMiddleware returns a middleware that checks for the presence of a specific cookie.
func MustCookieMiddleware(l logger.Logger, cookieName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := r.Cookie(cookieName)
			if err != nil {
				// Extract context, bind logger with request ID
				ctx := common.GetContext(r)
				log := logger.LoggerWithKey(l, ctx, common.ContextKeyRequestID)

				log.Warn("Missing required cookie",
					log.ToString("cookie_n``ame", cookieName),
				)

				details := "cookie '" + cookieName + "' not found"
				// Respond with error if the cookie is missing
				handlers.RespondWithError(log, w, ErrorMissingCookie, details)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
