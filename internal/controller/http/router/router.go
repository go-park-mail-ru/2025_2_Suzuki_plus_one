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
	r.Get("/media/my", h.GetMediaMy)
	r.Get("/media/recommendations", h.GetMediaRecommendations)
	r.Get(fmt.Sprintf("/media/{%s}", handlers.PathParamGetMediaID), h.GetMedia)
	r.Get(fmt.Sprintf("/media/{%s}/actor", handlers.PathParamGetMediaActorID), h.GetMediaActor)
	r.Get(fmt.Sprintf("/media/{%s}/like", handlers.PathParamGetMediaLikeID), h.GetMediaLike)
	r.Put(fmt.Sprintf("/media/{%s}/like", handlers.PathParamPutMediaLikeID), h.PutMediaLike)
	r.Delete(fmt.Sprintf("/media/{%s}/like", handlers.PathParamDeleteMediaLikeID), h.DeleteMediaLike)

	// Genre
	// r.Get("/genre/all", h.GetGenreAll)
	
	
	r.Get(fmt.Sprintf("/actor/{%s}", handlers.PathParamGetActorID), h.GetActor)
	r.Get(fmt.Sprintf("/actor/{%s}/media", handlers.PathParamGetActorMediaID), h.GetActorMedia)
	r.Get("/search", h.GetSearch)

	// Object
	r.Get("/object", h.GetObjectMedia)
	r.Get("/media/watch", h.GetMediaWatch)

	// Auth
	r.Group(func(r chi.Router) {
		// Must contain refresh token cookie
		r.Use(middleware.MustCookieMiddleware(l, handlers.CookieRefreshTokenName))
		r.Get("/auth/refresh", h.GetAuthRefresh)
	})

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
		r.Post("/user/me/update/avatar", h.PostUserMeUpdateAvatar)
		r.Post("/user/me/update/password", h.PostUserMeUpdatePassword)
	})

	// Appeal
	r.Get("/appeal/all", h.GetAppealAll)
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Authenticator(common.TokenAuth))
		r.Get("/appeal/my", h.GetAppealMy)
		r.Post("/appeal/new", h.PostAppealNew)
		r.Get(fmt.Sprintf("/appeal/{%s}", handlers.PathParamGetAppealID), h.GetAppeal)
		r.Put(fmt.Sprintf("/appeal/{%s}/resolve", handlers.PathParamGetAppealID), h.PutAppealResolve)
		r.Post(fmt.Sprintf("/appeal/{%s}/message", handlers.PathParamGetAppealID), h.PostAppealMessage)
		r.Get(fmt.Sprintf("/appeal/{%s}/message", handlers.PathParamGetAppealID), h.GetAppealMessage)
	})

	return r
}
