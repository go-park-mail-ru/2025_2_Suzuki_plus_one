package handlers

import (
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type Handlers struct {
	logger logger.Logger

	GetMovieRecommendationsUseCase controller.GetMovieRecommendationsUsecase
}

func NewHandlers(
	getMovieRecommendationsUsecase controller.GetMovieRecommendationsUsecase,
	logger logger.Logger,
) *Handlers {
	return &Handlers{
		logger:                         logger,
		GetMovieRecommendationsUseCase: getMovieRecommendationsUsecase,
	}
}
