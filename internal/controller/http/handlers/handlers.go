package handlers

import (
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type Handlers struct {
	logger logger.Logger

	GetMoviesUseCase controller.GetMoviesUsecase
}

func NewHandlers(getMoviesUsecase controller.GetMoviesUsecase, logger logger.Logger) *Handlers {
	return &Handlers{
		logger: logger,
		GetMoviesUseCase: getMoviesUsecase,
	}
}
