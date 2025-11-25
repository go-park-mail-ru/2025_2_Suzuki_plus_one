package http

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type GetUserMeUseCase struct {
	logger     logger.Logger
	userRepo   UserRepository
	objectRepo ObjectRepository
}

func NewGetUserMeUseCase(
	logger logger.Logger,
	userRepo UserRepository,
	objectRepo ObjectRepository,
) *GetUserMeUseCase {
	return &GetUserMeUseCase{
		logger:     logger,
		userRepo:   userRepo,
		objectRepo: objectRepo,
	}
}

func (uc *GetUserMeUseCase) Execute(ctx context.Context, input dto.GetUserMeInput) (dto.GetUserMeOutput, *dto.Error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContextKeyRequestID)

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_user_me",
			entity.ErrGetUserMeParamsInvalid,
			err.Error(),
		)
		return dto.GetUserMeOutput{}, &derr
	}

	// Compare user ID from session with requested user ID
	userID, err := common.ValidateToken(input.AccessToken)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_user_me",
			entity.ErrGetUserMeSessionNotFound,
			"access token is invalid",
		)
		log.Error("Failed to validate access token", log.ToError(err))
		return dto.GetUserMeOutput{}, &derr
	}

	// Get user info
	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_user_me",
			entity.ErrUserNotFound,
			err.Error(),
		)
		log.Error("Failed to get user by ID", log.ToError(err))
		return dto.GetUserMeOutput{}, &derr
	}

	log.Debug("Fetched user info", "user", log.ToAny("user", user))

	output := dto.GetUserMeOutput{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DateOfBirth: dto.JSONDate{Time: user.DateOfBirth},
		PhoneNumber: user.PhoneNumber,
	}
	// Get s3 key for user avatar
	avatarKey, err := uc.userRepo.GetUserAvatarKey(ctx, user.ID)
	if err != nil {
		log.Error("Failed to get presigned URL for user avatar", log.ToError(err))
	} else {
		// Generate public s3 URL for avatar
		avatarObject, err := uc.objectRepo.GeneratePublicURL(ctx, avatarKey.BucketName, avatarKey.Key)
		if err != nil {
			log.Error("Failed to get public URL for user avatar", log.ToError(err))
		} else {
			output.AvatarURL = avatarObject.URL
			log.Debug("Found s3 url for user", "url", output.AvatarURL, "userID", user.ID)
		}
	}

	// Assume that avatar is blank if it is not found
	// if output.AvatarURL == "" {
	// 	output.AvatarURL = "can't find avatar :("
	// }
	return output, nil
}
