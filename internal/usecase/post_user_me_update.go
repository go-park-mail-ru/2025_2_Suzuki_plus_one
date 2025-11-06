package usecase

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
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContexKeyRequestID)

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/post_user_me_update",
			err,
			"Invalid post user me update input parameters",
		)
		return dto.PostUserMeUpdateOutput{}, &derr
	}

	currentUser, derr := uc.getUserMeUseCase.Execute(ctx, dto.GetUserMeInput{
		AccessToken: input.AccessToken,
	})
	if derr != nil {
		log.Error(
			"Error extracting user with getUserMeUseCase",
			log.ToAny("derr", derr),
		)
		return dto.PostUserMeUpdateOutput{}, derr
	}

	// Update user in repository
	updatedUser, err := uc.userRepo.UpdateUser(
		ctx,
		currentUser.ID,
		input.Username,
		input.Email,
		input.DateOfBirth.GoString(),
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

	// Prepare output DTO
	output := dto.PostUserMeUpdateOutput{}
	output.ID = updatedUser.ID
	output.Username = updatedUser.Username
	output.Email = updatedUser.Email
	output.PhoneNumber = updatedUser.PhoneNumber
	output.DateOfBirth = input.DateOfBirth
	// Keep the existing avatar URL
	output.AvatarURL = currentUser.AvatarURL

	return output, nil
}
