package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
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
			"usecase/get_auth_signout",
			entity.ErrGetAuthSignOutInvalidParams,
			err.Error(),
		)
		return dto.GetAppealAllOutput{}, &derr
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
		return dto.GetAppealAllOutput{}, &derr
	}
	log.ToInt("userID: ", int(userID))

	//must be changed later
	return nil, nil
}
