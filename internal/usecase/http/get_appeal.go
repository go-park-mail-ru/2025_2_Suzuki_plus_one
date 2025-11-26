package http

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type GetAppealUseCase struct {
	logger     logger.Logger
	appealRepo AppealRepository
}

func NewGetAppealUseCase(
	logger logger.Logger,
	appealRepo AppealRepository,
) *GetAppealUseCase {
	return &GetAppealUseCase{
		logger:     logger,
		appealRepo: appealRepo,
	}
}

func (uc *GetAppealUseCase) Execute(ctx context.Context, input dto.GetAppealInput) (dto.GetAppealOutput, *dto.Error) {
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContextKeyRequestID)
	log.Debug("GetAppealUseCase called")

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_appeal",
			entity.ErrGetAppealFailed,
			err.Error(),
		)
		return dto.GetAppealOutput{}, &derr
	}

	// Validate JWT token, get user ID
	_, err := common.ValidateToken(input.AccessToken)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_appeal",
			entity.ErrGetAppealFailed,
			err.Error(),
		)
		log.Error("Invalid access token", log.ToError(err))
		return dto.GetAppealOutput{}, &derr
	}

	// Get appeal ids from repository
	appealByID, err := uc.appealRepo.GetAppealByID(ctx, input.AppealId)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_appeal",
			entity.ErrGetAppealFailed,
			err.Error(),
		)
		log.Error("Failed to get appeal ids", log.ToError(err))
		return dto.GetAppealOutput{}, &derr
	}

	output := dto.GetAppealOutput{
		Appeal: dto.Appeal{
			Appeal:    *appealByID,
			CreatedAt: dto.NewJSONDateTime(appealByID.CreatedAt),
			UpdatedAt: dto.NewJSONDateTime(appealByID.UpdatedAt),
		},
	}

	log.Debug("Successfully got details for a specific appeal", "AppealID", log.ToAny("AppealID", input.AppealId))
	return output, nil
}
