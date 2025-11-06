package http

import (
	"fmt"
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
	r.Use(chiMiddleware.RequestID)
	r.Use(middleware.GetLogging(l))
	r.Use(middleware.SetCors(origin))
	r.Use(middleware.SetJSON)
	r.Use(chiMiddleware.CleanPath)
	r.Use(chiMiddleware.StripSlashes)
	r.Use(chiMiddleware.Recoverer)
	r.Use(jwtauth.Verifier(common.TokenAuth))

	// Follow swagger order

	// Content
	r.Get("/movie/recommendations", h.GetMovieRecommendations)
	r.Get(fmt.Sprintf("/media/{%s}", handlers.PathParamGetMediaID), h.GetMedia)
	r.Get(fmt.Sprintf("/actor/{%s}", handlers.PathParamGetActorID), h.GetActor)

	// Object
	r.Get("/object", h.GetObjectMedia)
	r.Get("/media/watch", h.GetMediaWatch)

	// Auth
	r.Get("/auth/refresh", h.GetAuthRefresh)
	r.Post("/auth/signin", h.PostAuthSignIn)
	r.Post("/auth/signup", h.PostAuthSignUp)
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Authenticator(common.TokenAuth))
		r.Get("/auth/signout", h.GetAuthSignOut)
	})

	// User
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Authenticator(common.TokenAuth))
		r.Get("/user/me", h.GetUserMe)
		r.Post("/user/me/update", h.PostUserMeUpdate)
	})

	return r
}
