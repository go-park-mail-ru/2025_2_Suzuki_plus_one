package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
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
	r.Use(jwtauth.Verifier(common.TokenAuth))

	// Handlers
	r.Get("/movie/recommendations", h.GetMovieRecommendations)
	r.Get("/object", h.GetObjectMedia)

	// Auth routes
	r.Post("/auth/signin", h.PostAuthSignIn)
	r.Get("/auth/refresh", h.GetAuthRefresh)
	r.Post("/auth/signup", h.PostAuthSignUp)

	return r
}
