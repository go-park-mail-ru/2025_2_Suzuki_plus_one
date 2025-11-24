package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type PutAppealResolveUseCase struct {
	logger     logger.Logger
	appealRepo AppealRepository
}

func NewPutAppealResolveUseCase(
	logger logger.Logger,
	appealRepo AppealRepository,
) *PutAppealResolveUseCase {
	return &PutAppealResolveUseCase{
		logger:     logger,
		appealRepo: appealRepo,
	}
}

func (uc *PutAppealResolveUseCase) Execute(ctx context.Context, input dto.PutAppealResolveInput) (dto.PutAppealResolveOutput, *dto.Error) {
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContextKeyRequestID)
	log.Debug("PutAppealResolveUseCase called")

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/put_appeal_resolve",
			entity.ErrPutAppealResolve,
			err.Error(),
		)
		return dto.PutAppealResolveOutput{}, &derr
	}

	// Validate JWT token, get user ID
	_, err := common.ValidateToken(input.AccessToken)
	if err != nil {
		derr := dto.NewError(
			"usecase/put_appeal_resolve",
			entity.ErrWrongAccessToken,
			err.Error(),
		)
		log.Error("Invalid access token", log.ToError(err))
		return dto.PutAppealResolveOutput{}, &derr
	}

	err = uc.appealRepo.UpdateAppealStatus(ctx, input.AppealId, "resolved")
	if err != nil {
		derr := dto.NewError(
			"usecase/put_appeal_resolve",
			entity.ErrPutAppealResolve,
			err.Error(),
		)
		log.Error("Couldn't resolve appeal", log.ToError(err))
		return dto.PutAppealResolveOutput{}, &derr
	}

	output := dto.PutAppealResolveOutput{
		ID:      input.AppealId,
		Message: "Appeal resolved successfully.",
	}

	log.Debug("Successfully resolved an appeal", "appealID", log.ToAny("appealID", input.AppealId))
	return output, nil
}
