package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

func NewServer(router http.Handler) http.Handler {
	// Create server
	server := chi.NewMux()
	// server.Handle("/", router) // for http.ServeMux
	server.Mount("/", router)
	return server
}

func StartServer(router http.Handler, serveString string, l logger.Logger) {
	server := NewServer(router)

	// Start the server
	if err := http.ListenAndServe(serveString, server); err != nil {
		l.Fatal("Server not started 'cause of error", l.ToError(err))
	}
	l.Info("Server started serving")
}
