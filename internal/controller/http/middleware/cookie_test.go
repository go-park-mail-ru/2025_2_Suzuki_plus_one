package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestMustCookieMiddleware(t *testing.T) {
	testLogger := logger.NewZapLogger(true)

	tests := []struct {
		name           string
		cookieName     string
		cookies        []*http.Cookie
		expectedStatus int
		shouldCallNext bool
	}{
		{
			name:       "Cookie exists - should call next handler",
			cookieName: "session",
			cookies: []*http.Cookie{
				{Name: "session", Value: "abc123"},
			},
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
		},
		{
			name:           "Cookie missing - should not call next handler",
			cookieName:     "session",
			cookies:        []*http.Cookie{},
			expectedStatus: http.StatusUnauthorized,
			shouldCallNext: false,
		},
		{
			name:       "Different cookie exists but required cookie missing",
			cookieName: "session",
			cookies: []*http.Cookie{
				{Name: "other", Value: "value"},
			},
			expectedStatus: http.StatusUnauthorized,
			shouldCallNext: false,
		},
		{
			name:       "Multiple cookies including required one",
			cookieName: "auth",
			cookies: []*http.Cookie{
				{Name: "session", Value: "abc123"},
				{Name: "auth", Value: "token123"},
				{Name: "preferences", Value: "dark"},
			},
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Флаг для проверки вызова следующего обработчика
			nextCalled := false

			// Создаем тестовый обработчик
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			})

			// Создаем middleware
			middleware := MustCookieMiddleware(testLogger, tt.cookieName)
			handler := middleware(nextHandler)

			// Создаем запрос
			req := httptest.NewRequest("GET", "/test", nil)

			// Добавляем cookies к запросу
			for _, cookie := range tt.cookies {
				req.AddCookie(cookie)
			}

			// Добавляем контекст с request ID
			ctx := context.WithValue(req.Context(), common.ContextKeyRequestID, "test-request-id")
			req = req.WithContext(ctx)

			// Создаем ResponseRecorder
			rr := httptest.NewRecorder()

			// Вызываем обработчик
			handler.ServeHTTP(rr, req)

			// Проверяем статус код
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Проверяем, был ли вызван следующий обработчик
			assert.Equal(t, tt.shouldCallNext, nextCalled)

			// Проверяем тело ответа для случая ошибки
			if !tt.shouldCallNext {
				assert.Contains(t, rr.Body.String(), "Missing required cookie")
			} else {
				assert.Equal(t, "OK", rr.Body.String())
			}
		})
	}
}
