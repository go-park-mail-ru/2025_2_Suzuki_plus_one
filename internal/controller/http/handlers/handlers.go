package handlers

import (
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type Handlers struct {
	Logger logger.Logger

	GetMovieRecommendationsUseCase controller.GetMovieRecommendationsUsecase
	GetObjectMediaUseCase          controller.GetObjectUsecase
	PostAuthSignInUseCase          controller.PostAuthSignInUsecase
	GetAuthRefreshUseCase          controller.GetAuthRefreshUsecase
}

func NewHandlers(
	logger logger.Logger,
	getMovieRecommendationsUseCase controller.GetMovieRecommendationsUsecase,
	getObjectMediaUseCase controller.GetObjectUsecase,
	postAuthSignInUseCase controller.PostAuthSignInUsecase,
	getAuthRefreshUseCase controller.GetAuthRefreshUsecase,
) *Handlers {
	return &Handlers{
		Logger:                         logger,
		GetMovieRecommendationsUseCase: getMovieRecommendationsUseCase,
		GetObjectMediaUseCase:          getObjectMediaUseCase,
		PostAuthSignInUseCase:          postAuthSignInUseCase,
		GetAuthRefreshUseCase:          getAuthRefreshUseCase,
	}
}
