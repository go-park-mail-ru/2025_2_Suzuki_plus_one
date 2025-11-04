package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type GetUserMeUseCase struct {
	logger      logger.Logger
	userRepo    UserRepository
	sessionRepo SessionRepository
	objectRepo  ObjectRepository
}

func NewGetUserMeUseCase(
	logger logger.Logger,
	userRepo UserRepository,
	sessionRepo SessionRepository,
	objectRepo ObjectRepository,
) *GetUserMeUseCase {
	return &GetUserMeUseCase{
		logger:      logger,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		objectRepo:  objectRepo,
	}
}

func (uc *GetUserMeUseCase) Execute(ctx context.Context, input dto.GetUserMeInput) (dto.GetUserMeOutput, *dto.Error) {
	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_user_me",
			entity.ErrGetUserMeParamsInvalid,
			err.Error(),
		)
		return dto.GetUserMeOutput{}, &derr
	}

	// Get session by access token
	userID, err := uc.sessionRepo.GetUserIDByToken(ctx, input.AccessToken)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_user_me",
			entity.ErrGetUserMeSessionNotFound,
			err.Error(),
		)
		uc.logger.Error("Failed to get user ID by token", uc.logger.ToError(err))
		return dto.GetUserMeOutput{}, &derr
	}

	// Compare user ID from session with requested user ID
	userIDToken, err := common.ValidateToken(input.AccessToken)
	if err != nil || userID != userIDToken {
		derr := dto.NewError(
			"usecase/get_user_me",
			entity.ErrGetUserMeSessionNotFound,
			"access token is invalid",
		)
		uc.logger.Error("Failed to validate access token", uc.logger.ToError(err))
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
		uc.logger.Error("Failed to get user by ID", uc.logger.ToError(err))
		return dto.GetUserMeOutput{}, &derr
	}

	var dateOfBirth time.Time
	if user.DateOfBirth != "" {
		parsedDOB, err := time.Parse("2006-01-02", user.DateOfBirth)
		if err != nil {
			uc.logger.Error("Failed to parse user date of birth", uc.logger.ToError(err))
		} else {
			dateOfBirth = parsedDOB
		}
	}

	output := dto.GetUserMeOutput{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DateOfBirth: dateOfBirth,
		PhoneNumber: user.PhoneNumber,
	}
	// Get s3 key for user avatar
	avatarKey, err := uc.userRepo.GetUserAvatarKey(ctx, user.ID)
	if err != nil {
		uc.logger.Error("Failed to get presigned URL for user avatar", uc.logger.ToError(err))
	} else {
		// NOTE: This is fucking awful, please refactor later
		// Split avatarKey into bucket and object name
		// Trim first slash if exists
		if len(avatarKey) > 0 && avatarKey[0] == '/' {
			avatarKey = avatarKey[1:]
		}
		parts := strings.SplitN(avatarKey, "/", 2)
		var bucket, key string
		if len(parts) == 2 {
			bucket = parts[0]
			key = parts[1]
		} else {
			bucket = avatarKey
			key = ""
		}

		// Generate presigned URL for avatar
		avatarObject, err := uc.objectRepo.GetPublicObject(ctx, key, bucket)
		if err != nil {
			uc.logger.Error("Failed to get presigned URL for user avatar", uc.logger.ToError(err))
		} else {
			output.AvatarURL = avatarObject.URL
			uc.logger.Debug("Found s3 url for user", "url", output.AvatarURL, "userID", user.ID)
		}
	}
	return output, nil

}
