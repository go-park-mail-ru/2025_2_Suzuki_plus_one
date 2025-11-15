package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type GetAppealAllUseCase struct {
	logger     logger.Logger
	appealRepo AppealRepository
}

func NewGetAppealAllUseCase(
	logger logger.Logger,
	appealRepo AppealRepository,
) *GetAppealAllUseCase {
	return &GetAppealAllUseCase{
		logger:     logger,
		appealRepo: appealRepo,
	}
}

func (uc *GetAppealAllUseCase) Execute(ctx context.Context, input dto.GetAppealAllInput) (dto.GetAppealAllOutput, *dto.Error) {
	log := logger.LoggerWithKey(uc.logger, ctx, "request_id")
	log.Debug("GetAppealAllUseCase called")

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_appeal_all",
			entity.ErrGetAuthSignOutInvalidParams,
			err.Error(),
		)
		return dto.GetAppealAllOutput{}, &derr
	}

	var tag *string
	if input.Tag != "" {
		tag = &input.Tag
	}
	var status *string
	if input.Status != "" {
		status = &input.Status
	}

	// We exxpect to fail on validation without limit and offset

	// Get appeal IDs
	appealIDs, err := uc.appealRepo.GetAppealIDsAll(ctx, tag, status, input.Limit, input.Offset)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_appeal_all",
			entity.ErrGetAppealAllFailed,
			"AppealRepository.GetAppealIDsAll failed: "+err.Error(),
		)
		return dto.GetAppealAllOutput{}, &derr
	}

	// Get full appeal info
	appeals := make([]dto.GetAppealOutput, 0, len(appealIDs))
	for _, appealID := range appealIDs {
		appeal, err := uc.appealRepo.GetAppealByID(ctx, appealID)
		if err != nil {
			derr := dto.NewError(
				"usecase/get_appeal_all",
				entity.ErrGetAppealAllFailed,
				"AppealRepository.GetAppealByID failed: "+err.Error(),
			)
			return dto.GetAppealAllOutput{}, &derr
		}
		appeals = append(appeals, dto.GetAppealOutput{
			Appeal: dto.Appeal{
				Appeal:    *appeal,
				CreatedAt: dto.NewJSONDateTime(appeal.CreatedAt),
				UpdatedAt: dto.NewJSONDateTime(appeal.UpdatedAt),
			},
		})
	}

	log.Debug("GetAppealAllUseCase completed successfully")
	return dto.GetAppealAllOutput{
		Appeals: appeals,
	}, nil
}
