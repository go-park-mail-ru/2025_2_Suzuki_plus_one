package usecase

import "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"

type GetAppealMyUseCase struct {
	logger     logger.Logger
	appealRepo AppealRepository
}

func NewGetAppealMyUseCase(
	logger logger.Logger,
	appealRepo AppealRepository,
) *GetAppealMyUseCase {
	return &GetAppealMyUseCase{
		logger:     logger,
		appealRepo: appealRepo,
	}
}
