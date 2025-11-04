package handlers

import (
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type Handlers struct {
	logger logger.Logger

	GetMovieRecommendationsUseCase controller.GetMovieRecommendationsUsecase
	GetObjectMediaUseCase          controller.GetObjectUsecase
}

func NewHandlers(
	logger logger.Logger,
	getMovieRecommendationsUsecase controller.GetMovieRecommendationsUsecase,
	getObjectMediaUsecase controller.GetObjectUsecase,
) *Handlers {
	return &Handlers{
		logger:                         logger,
		GetMovieRecommendationsUseCase: getMovieRecommendationsUsecase,
		GetObjectMediaUseCase:          getObjectMediaUsecase,
	}
}
