package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type PostAppealNewUseCase struct {
	logger     logger.Logger
	appealRepo AppealRepository
}

func NewPostAppealNewUseCase(
	logger logger.Logger,
	appealRepo AppealRepository,
) *PostAppealNewUseCase {
	return &PostAppealNewUseCase{
		logger:     logger,
		appealRepo: appealRepo,
	}
}

func zipStringToTen(toZip string) string {
	if len(toZip) <= 10 {
		return toZip
	}

	return toZip[0:10]
}

func (uc *PostAppealNewUseCase) Execute(ctx context.Context, input dto.PostAppealNewInput) (dto.PostAppealNewOutput, *dto.Error) {
	log := logger.LoggerWithKey(uc.logger, ctx, "request_id")
	log.Debug("PostAppealNewUseCase called")

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/post_appeal_new",
			entity.ErrPostAppealNew,
			err.Error(),
		)
		return dto.PostAppealNewOutput{}, &derr
	}

	// Validate JWT token, get user ID
	userID, err := common.ValidateToken(input.AccessToken)
	if err != nil {
		derr := dto.NewError(
			"usecase/post_appeal_new",
			entity.ErrWrongAccessToken,
			err.Error(),
		)
		log.Error("Invalid access token", log.ToError(err))
		return dto.PostAppealNewOutput{}, &derr
	}

	appealID, err := uc.appealRepo.CreateAppeal(ctx, userID, input.Tag, zipStringToTen(input.Message))
	if err != nil {
		derr := dto.NewError(
			"usecase/post_appeal_new",
			entity.ErrPostAppealNew,
			err.Error(),
		)
		log.Error("Couldn't create appeal", log.ToError(err))
		return dto.PostAppealNewOutput{}, &derr
	}

	output := dto.PostAppealNewOutput{
		ID: appealID,
	}

	log.Debug("Successfully created new appeal", "appealID", log.ToAny("appealID", appealID))
	return output, nil
}
