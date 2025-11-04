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
}
