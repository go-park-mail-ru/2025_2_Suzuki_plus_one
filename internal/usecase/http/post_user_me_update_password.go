package http

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type PostUserMeUpdatePasswordUseCase struct {
	logger      logger.Logger
	userRepo    UserRepository
	sessionRepo SessionRepository
}

func NewPostUserMeUpdatePasswordUseCase(
	logger logger.Logger,
	userRepo UserRepository,
	sessionRepo SessionRepository,
) *PostUserMeUpdatePasswordUseCase {
	return &PostUserMeUpdatePasswordUseCase{
		logger:      logger,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (uc *PostUserMeUpdatePasswordUseCase) Execute(ctx context.Context, input dto.PostUserMeUpdatePasswordInput) (dto.PostUserMeUpdatePasswordOutput, *dto.Error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContextKeyRequestID)

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/post_user_me_update_password",
			err,
			"Invalid post user me update input parameters",
		)
		return dto.PostUserMeUpdatePasswordOutput{}, &derr
	}

	// Get user ID from cache using access token
	userID, err := uc.sessionRepo.GetUserIDByAccessToken(ctx, input.AccessToken)
	if err != nil {
		derr := dto.NewError(
			"usecase/post_user_me_update_password/get_user_id_by_access_token",
			err,
			"Failed to get user ID by access token",
		)
		log.Error(
			"Error getting user ID by access token",
			log.ToError(err),
		)
		return dto.PostUserMeUpdatePasswordOutput{}, &derr
	}

	// Get current user data
	currentUser, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		derr := dto.NewError(
			"usecase/post_user_me_update_password/get_user_by_id",
			err,
			"Failed to get current user data",
		)
		log.Error(
			"Error getting current user data",
			log.ToError(err),
		)
		return dto.PostUserMeUpdatePasswordOutput{}, &derr
	}

	// Validate current password
	err = common.ValidateHashedPasswordBcrypt(currentUser.PasswordHash, input.CurrentPassword)
	if err != nil {
		derr := dto.NewError(
			"usecase/post_user_me_update_password/validate_current_password",
			entity.ErrPostUserMeUpdatePasswordCurrentPasswordMismatch,
			err.Error(),
		)
		log.Error(
			"Current password does not match",
			log.ToAny("current password by user request", input.CurrentPassword),
			log.ToAny("stored password hash", currentUser.PasswordHash),
		)
		return dto.PostUserMeUpdatePasswordOutput{}, &derr
	}

	// Update user password in repository
	err = uc.userRepo.UpdateUserPassword(
		ctx,
		userID,
		input.NewPassword,
	)
	if err != nil {
		derr := dto.NewError(
			"usecase/post_user_me_update_password/update_user_password",
			err,
			"Failed to update user password",
		)
		log.Error(
			"Error updating user password",
			log.ToError(err),
		)
		return dto.PostUserMeUpdatePasswordOutput{}, &derr
	}

	output := dto.PostUserMeUpdatePasswordOutput{}

	return output, nil
}
