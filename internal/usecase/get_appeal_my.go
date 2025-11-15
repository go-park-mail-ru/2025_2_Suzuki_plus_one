package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

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

func (uc *GetAppealMyUseCase) Execute(ctx context.Context, input dto.GetAppealMyInput) (dto.GetAppealMyOutput, *dto.Error) {
	log := logger.LoggerWithKey(uc.logger, ctx, "request_id")
	log.Debug("GetAppealMyUseCase called")

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_appeal_my",
			entity.ErrGetAppealMyFailed,
			err.Error(),
		)
		return dto.GetAppealMyOutput{}, &derr
	}

	// Validate JWT token, get user ID
	userID, err := common.ValidateToken(input.AccessToken)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_appeal_my",
			entity.ErrWrongAccessToken,
			err.Error(),
		)
		log.Error("Invalid access token", log.ToError(err))
		return dto.GetAppealMyOutput{}, &derr
	}

	// Get appeal ids from repository
	appealIDs, err := uc.appealRepo.GetAppealIDsByUserID(ctx, userID)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_appeal_my",
			entity.ErrGetAppealMyFailed,
			err.Error(),
		)
		log.Error("Failed to get appeal ids", log.ToError(err))
		return dto.GetAppealMyOutput{}, &derr
	}

	output := dto.GetAppealMyOutput{
		Appeals: make([]dto.GetAppealOutput, 0, len(appealIDs)),
	}

	// Get appeals by ids
	for _, appealID := range appealIDs {
		appealEntity, err := uc.appealRepo.GetAppealByID(ctx, appealID)
		if err != nil {
			derr := dto.NewError(
				"usecase/get_appeal_my",
				entity.ErrGetAppealMyFailed,
				err.Error(),
			)
			log.Error("Failed to get appeal by ID", log.ToError(err))
			return dto.GetAppealMyOutput{}, &derr
		}

		appealDTO := dto.GetAppealOutput{
			Appeal: *appealEntity,
		}
		output.Appeals = append(output.Appeals, appealDTO)
	}

	log.Debug("Successfully fetched appeals for user", "user_id", log.ToAny("user_id", userID))
	return output, nil
}
