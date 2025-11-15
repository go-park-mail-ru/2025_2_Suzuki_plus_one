package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type GetAppealMessageUseCase struct {
	logger     logger.Logger
	appealRepo AppealRepository
}

func NewGetAppealMessageUseCase(
	logger logger.Logger,
	appealRepo AppealRepository,
) *GetAppealMessageUseCase {
	return &GetAppealMessageUseCase{
		logger:     logger,
		appealRepo: appealRepo,
	}
}

func (uc *GetAppealMessageUseCase) Execute(ctx context.Context, input dto.GetAppealMessageInput) (dto.GetAppealMessageOutput, *dto.Error) {
	log := logger.LoggerWithKey(uc.logger, ctx, "request_id")
	log.Debug("GetAppealMessageUseCase called")

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_appeal_message",
			entity.ErrGetAppealMessageInvalidParams,
			err.Error(),
		)
		return dto.GetAppealMessageOutput{}, &derr
	}

	// Validate JWT token, get user ID
	_, err := common.ValidateToken(input.AccessToken)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_appeal_message",
			entity.ErrWrongAccessToken,
			err.Error(),
		)
		log.Error("Invalid access token", log.ToError(err))
		return dto.GetAppealMessageOutput{}, &derr
	}

	// Get appeal message in repository
	messages, err := uc.appealRepo.GetAppealMessagesByID(ctx, input.AppealID)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_appeal_message",
			entity.ErrGetAppealMessageInvalidParams,
			"AppealRepository.GetAppealMessagesByID failed: "+err.Error(),
		)
		return dto.GetAppealMessageOutput{}, &derr
	}

	output := dto.GetAppealMessageOutput{
		Messages: make([]dto.AppealMessage, 0, len(messages)),
	}
	for _, msg := range messages {
		output.Messages = append(output.Messages, dto.AppealMessage{
			AppealMessage: msg,
			CreatedAt:     dto.NewJSONDateTime(msg.CreatedAt),
		})
	}

	return output, nil
}
