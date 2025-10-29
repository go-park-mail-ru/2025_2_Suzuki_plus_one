package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http/handlers"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http/middleware"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

func InitRouter(h *handlers.Handlers, l logger.Logger, origin string) http.Handler {
	r := chi.NewRouter()

	// Middlewares, in the correct order
	r.Use(middleware.GetLogging(l)) // Logger have to be first
	r.Use(middleware.SetCors(origin))
	r.Use(middleware.SetJSON)
	r.Use(chiMiddleware.CleanPath)
	r.Use(chiMiddleware.StripSlashes)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RequestID)

	// Handlers
	r.Get("/movies", h.GetMovies)

	return r
}
