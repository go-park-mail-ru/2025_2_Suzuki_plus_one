package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/db"
)

func NewMockServer() *Server {
	return &Server{
		db:     &db.DataBase{},
		server: http.NewServeMux(),
	}
}

func TestMovies(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/movies", nil)
	w := httptest.NewRecorder()

	NewMockServer().getAllMovies(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}
