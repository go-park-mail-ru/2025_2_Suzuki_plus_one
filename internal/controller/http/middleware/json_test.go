package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetJSON(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		existingHeader string
		expectedHeader string
		description    string
	}{
		{
			name:           "GET request without existing header",
			method:         "GET",
			existingHeader: "",
			expectedHeader: "application/json",
			description:    "Should set Content-Type to application/json",
		},
		{
			name:           "POST request without existing header",
			method:         "POST",
			existingHeader: "",
			expectedHeader: "application/json",
			description:    "Should set Content-Type for POST request",
		},
		{
			name:           "PUT request without existing header",
			method:         "PUT",
			existingHeader: "",
			expectedHeader: "application/json",
			description:    "Should set Content-Type for PUT request",
		},
		{
			name:           "DELETE request without existing header",
			method:         "DELETE",
			existingHeader: "",
			expectedHeader: "application/json",
			description:    "Should set Content-Type for DELETE request",
		},
		{
			name:           "Request with existing Content-Type header",
			method:         "GET",
			existingHeader: "text/html",
			expectedHeader: "application/json",
			description:    "Should override existing Content-Type header",
		},
		{
			name:           "OPTIONS request",
			method:         "OPTIONS",
			existingHeader: "",
			expectedHeader: "application/json",
			description:    "Should set Content-Type for OPTIONS request",
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
				w.Write([]byte(`{"status": "ok"}`))
			})

			// Создаем middleware
			handler := SetJSON(nextHandler)

			// Создаем запрос
			req := httptest.NewRequest(tt.method, "/test", nil)

			// Создаем ResponseRecorder
			rr := httptest.NewRecorder()

			// Устанавливаем существующий заголовок если нужно
			if tt.existingHeader != "" {
				rr.Header().Set("Content-Type", tt.existingHeader)
			}

			// Вызываем обработчик
			handler.ServeHTTP(rr, req)

			// Проверяем заголовок Content-Type
			assert.Equal(t, tt.expectedHeader, rr.Header().Get("Content-Type"), tt.description)

			// Проверяем, что следующий обработчик был вызван
			assert.True(t, nextCalled, "Next handler should be called")

			// Проверяем статус код
			assert.Equal(t, http.StatusOK, rr.Code)

			// Проверяем тело ответа
			assert.Equal(t, `{"status": "ok"}`, rr.Body.String())
		})
	}
}
