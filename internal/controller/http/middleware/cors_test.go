package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetCors(t *testing.T) {
	tests := []struct {
		name            string
		frontendOrigin  string
		method          string
		expectedOrigin  string
		expectedMethods string
		expectedHeaders string
		expectedStatus  int
		shouldCallNext  bool
		description     string
	}{
		{
			name:            "Regular request with specific origin",
			frontendOrigin:  "https://example.com",
			method:          "GET",
			expectedOrigin:  "https://example.com",
			expectedMethods: "GET, POST, PUT, DELETE, OPTIONS",
			expectedHeaders: "Content-Type, Authorization",
			expectedStatus:  http.StatusOK,
			shouldCallNext:  true,
			description:     "Should set CORS headers and call next handler",
		},
		{
			name:            "OPTIONS preflight request",
			frontendOrigin:  "https://example.com",
			method:          "OPTIONS",
			expectedOrigin:  "https://example.com",
			expectedMethods: "GET, POST, PUT, DELETE, OPTIONS",
			expectedHeaders: "Content-Type, Authorization",
			expectedStatus:  http.StatusOK,
			shouldCallNext:  false,
			description:     "Should handle OPTIONS request without calling next handler",
		},
		{
			name:            "POST request with wildcard origin",
			frontendOrigin:  "*",
			method:          "POST",
			expectedOrigin:  "*",
			expectedMethods: "GET, POST, PUT, DELETE, OPTIONS",
			expectedHeaders: "Content-Type, Authorization",
			expectedStatus:  http.StatusOK,
			shouldCallNext:  true,
			description:     "Should set wildcard origin for POST request",
		},
		{
			name:            "PUT request with localhost origin",
			frontendOrigin:  "http://localhost:3000",
			method:          "PUT",
			expectedOrigin:  "http://localhost:3000",
			expectedMethods: "GET, POST, PUT, DELETE, OPTIONS",
			expectedHeaders: "Content-Type, Authorization",
			expectedStatus:  http.StatusOK,
			shouldCallNext:  true,
			description:     "Should set localhost origin for PUT request",
		},
		{
			name:            "DELETE request",
			frontendOrigin:  "https://api.example.com",
			method:          "DELETE",
			expectedOrigin:  "https://api.example.com",
			expectedMethods: "GET, POST, PUT, DELETE, OPTIONS",
			expectedHeaders: "Content-Type, Authorization",
			expectedStatus:  http.StatusOK,
			shouldCallNext:  true,
			description:     "Should set CORS headers for DELETE request",
		},
		{
			name:            "Empty origin",
			frontendOrigin:  "",
			method:          "GET",
			expectedOrigin:  "",
			expectedMethods: "GET, POST, PUT, DELETE, OPTIONS",
			expectedHeaders: "Content-Type, Authorization",
			expectedStatus:  http.StatusOK,
			shouldCallNext:  true,
			description:     "Should work with empty origin",
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

			// Создаем CORS middleware
			corsMiddleware := SetCors(tt.frontendOrigin)
			handler := corsMiddleware(nextHandler)

			// Создаем запрос
			req := httptest.NewRequest(tt.method, "/test", nil)

			// Создаем ResponseRecorder
			rr := httptest.NewRecorder()

			// Вызываем обработчик
			handler.ServeHTTP(rr, req)

			// Проверяем CORS заголовки
			assert.Equal(t, tt.expectedOrigin, rr.Header().Get("Access-Control-Allow-Origin"))
			assert.Equal(t, tt.expectedMethods, rr.Header().Get("Access-Control-Allow-Methods"))
			assert.Equal(t, tt.expectedHeaders, rr.Header().Get("Access-Control-Allow-Headers"))

			// Проверяем статус код
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Проверяем, был ли вызван следующий обработчик
			assert.Equal(t, tt.shouldCallNext, nextCalled, tt.description)

			// Проверяем тело ответа для не-OPTIONS запросов
			if tt.method != "OPTIONS" && tt.shouldCallNext {
				assert.Equal(t, "OK", rr.Body.String())
			} else if tt.method == "OPTIONS" {
				assert.Equal(t, "", rr.Body.String(), "OPTIONS request should have empty body")
			}
		})
	}
}
