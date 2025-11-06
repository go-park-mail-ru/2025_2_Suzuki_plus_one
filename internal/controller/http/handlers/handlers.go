package handlers

import (
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type Handlers struct {
	Logger logger.Logger

	GetMovieRecommendationsUseCase controller.GetMovieRecommendationsUseCase
	GetObjectMediaUseCase          controller.GetObjectUseCase
	PostAuthSignInUseCase          controller.PostAuthSignInUseCase
	GetAuthRefreshUseCase          controller.GetAuthRefreshUseCase
	PostAuthSignUpUseCase          controller.PostAuthSignUpUseCase
	GetAuthSignOutUseCase          controller.GetAuthSignOutUseCase
	GetUserMeUseCase               controller.GetUserMeUseCase
	GetActorUseCase                controller.GetActorUseCase
	GetMediaUseCase                controller.GetMediaUseCase
	GetMediaWatchUseCase           controller.GetMediaWatchUseCase
	PostUserMeUpdateUseCase        controller.PostUserMeUpdateUseCase
}

func NewHandlers(
	logger logger.Logger,
	getMovieRecommendationsUseCase controller.GetMovieRecommendationsUseCase,
	getObjectMediaUseCase controller.GetObjectUseCase,
	postAuthSignInUseCase controller.PostAuthSignInUseCase,
	getAuthRefreshUseCase controller.GetAuthRefreshUseCase,
	postAuthSignupUseCase controller.PostAuthSignUpUseCase,
	getAuthSignOutUseCase controller.GetAuthSignOutUseCase,
	getUserMeUseCase controller.GetUserMeUseCase,
	GetActorUseCase controller.GetActorUseCase,
	getMediaUseCase controller.GetMediaUseCase,
	getMediaWatchUseCase controller.GetMediaWatchUseCase,
	PostUserMeUpdateUseCase controller.PostUserMeUpdateUseCase,
) *Handlers {
	return &Handlers{
		Logger:                         logger,
		GetMovieRecommendationsUseCase: getMovieRecommendationsUseCase,
		GetObjectMediaUseCase:          getObjectMediaUseCase,
		PostAuthSignInUseCase:          postAuthSignInUseCase,
		GetAuthRefreshUseCase:          getAuthRefreshUseCase,
		PostAuthSignUpUseCase:          postAuthSignupUseCase,
		GetAuthSignOutUseCase:          getAuthSignOutUseCase,
		GetUserMeUseCase:               getUserMeUseCase,
		GetActorUseCase:                GetActorUseCase,
		GetMediaUseCase:                getMediaUseCase,
		GetMediaWatchUseCase:           getMediaWatchUseCase,
		PostUserMeUpdateUseCase:        PostUserMeUpdateUseCase,
	}
}
