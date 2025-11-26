package http

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type PostAppealMessageUseCase struct {
	logger     logger.Logger
	appealRepo AppealRepository
}

func NewPostAppealMessageUseCase(
	logger logger.Logger,
	appealRepo AppealRepository,
) *PostAppealMessageUseCase {
	return &PostAppealMessageUseCase{
		logger:     logger,
		appealRepo: appealRepo,
	}
}

func (uc *PostAppealMessageUseCase) Execute(ctx context.Context, input dto.PostAppealMessageInput) (dto.PostAppealMessageOutput, *dto.Error) {
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContextKeyRequestID)
	log.Debug("PostAppealMessageUseCase called")

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/post_appeal_message",
			entity.ErrPostAppealMessageInvalidParams,
			err.Error(),
		)
		return dto.PostAppealMessageOutput{}, &derr
	}

	// Validate JWT token, get user ID
	_, err := common.ValidateToken(input.AccessToken)
	if err != nil {
		derr := dto.NewError(
			"usecase/post_appeal_message",
			entity.ErrWrongAccessToken,
			err.Error(),
		)
		log.Error("Invalid access token", log.ToError(err))
		return dto.PostAppealMessageOutput{}, &derr
	}

	// Create appeal message in repository
	isResponse := false // TODO: determine if message is a response from support staff

	msgID, err := uc.appealRepo.CreateAppealMessage(ctx, input.AppealID, isResponse, input.Message)
	if err != nil {
		derr := dto.NewError(
			"usecase/post_appeal_message",
			entity.ErrPostAppealMessageInvalidParams,
			"AppealRepository.CreateAppealMessage failed: "+err.Error(),
		)
		return dto.PostAppealMessageOutput{}, &derr
	}

	output := dto.PostAppealMessageOutput{
		ID: msgID,
	}

	return output, nil
}
