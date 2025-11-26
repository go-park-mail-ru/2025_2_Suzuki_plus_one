package http

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type PostUserMeUpdateUseCase struct {
	logger           logger.Logger
	userRepo         UserRepository
	getUserMeUseCase *GetUserMeUseCase
}

func NewPostUserMeUpdateUseCase(
	logger logger.Logger,
	userRepo UserRepository,
	getUserMeUseCase *GetUserMeUseCase,
) *PostUserMeUpdateUseCase {
	return &PostUserMeUpdateUseCase{
		logger:           logger,
		userRepo:         userRepo,
		getUserMeUseCase: getUserMeUseCase,
	}
}

func (uc *PostUserMeUpdateUseCase) Execute(ctx context.Context, input dto.PostUserMeUpdateInput) (dto.PostUserMeUpdateOutput, *dto.Error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContextKeyRequestID)

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/post_user_me_update",
			err,
			"Invalid post user me update input parameters",
		)
		return dto.PostUserMeUpdateOutput{}, &derr
	}

	// Get user ID from access token
	userID, err := common.ValidateToken(input.AccessToken)
	if err != nil {
		derr := dto.NewError(
			"usecase/post_user_me_update",
			err,
			"Access token is invalid",
		)
		log.Error("Failed to validate access token", log.ToError(err))
		return dto.PostUserMeUpdateOutput{}, &derr
	}

	// Update user in repository
	_, err = uc.userRepo.UpdateUser(
		ctx,
		userID,
		input.Username,
		input.Email,
		input.DateOfBirth.Time,
		input.PhoneNumber,
	)

	if err != nil {
		derr := dto.NewError(
			"usecase/post_user_me_update/update_user",
			err,
			"Failed to update user",
		)
		log.Error("Failed to update user", log.ToError(err))
		return dto.PostUserMeUpdateOutput{}, &derr
	}

	// Query database again to get updated user info
	// OPTIMIZATION: consider returning updated fields directly
	getUserMeUseCaseOutput, derr := uc.getUserMeUseCase.Execute(ctx, dto.GetUserMeInput{
		AccessToken: input.AccessToken,
	})
	if derr != nil {
		log.Error(
			"Error extracting user with getUserMeUseCase",
			log.ToAny("derr", derr),
		)
		return dto.PostUserMeUpdateOutput{}, derr
	}

	return dto.PostUserMeUpdateOutput{
		GetUserMeOutput: getUserMeUseCaseOutput,
	}, nil
}
