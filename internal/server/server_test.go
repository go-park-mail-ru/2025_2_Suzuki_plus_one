// This file contains tests that hit actual endpoints

package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/db"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/models"
)

func NewMockServer() *Server {
	return &Server{
		db:     db.NewMockDB(),
		server: http.NewServeMux(),
	}
}

func TestMovies(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/movies", nil)
	w := httptest.NewRecorder()

	// Connect JSON middleware
	srv := NewMockServer()
	handler := forceJSONMiddleware(http.HandlerFunc(srv.getAllMovies))
	handler.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got [%s]", ct)
	}

	var movies []models.Movie
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&movies); err != nil {
		t.Errorf("failed to decode response body: %v", err)
	}
}

// TODO: Add more tests for other handlers and auth middleware at least
