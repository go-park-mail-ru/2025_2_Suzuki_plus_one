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
}

func NewHandlers(
	logger logger.Logger,
	getMovieRecommendationsUseCase controller.GetMovieRecommendationsUseCase,
	getObjectMediaUseCase controller.GetObjectUseCase,
	postAuthSignInUseCase controller.PostAuthSignInUseCase,
	getAuthRefreshUseCase controller.GetAuthRefreshUseCase,
	postAuthSignupUseCase controller.PostAuthSignUpUseCase,
) *Handlers {
	return &Handlers{
		Logger:                         logger,
		GetMovieRecommendationsUseCase: getMovieRecommendationsUseCase,
		GetObjectMediaUseCase:          getObjectMediaUseCase,
		PostAuthSignInUseCase:          postAuthSignInUseCase,
		GetAuthRefreshUseCase:          getAuthRefreshUseCase,
		PostAuthSignUpUseCase:          postAuthSignupUseCase,
	}
}
