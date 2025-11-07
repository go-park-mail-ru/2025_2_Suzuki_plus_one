package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type PostUserMeUpdateAvatarUseCase struct {
	logger      logger.Logger
	userRepo    UserRepository
	sessionRepo SessionRepository
	objectRepo  ObjectRepository
	assetRepo   AssetRepository
}

func NewPostUserMeUpdateAvatarUseCase(
	logger logger.Logger,
	userRepo UserRepository,
	sessionRepo SessionRepository,
	objectRepo ObjectRepository,
	assetRepo AssetRepository,
) *PostUserMeUpdateAvatarUseCase {
	return &PostUserMeUpdateAvatarUseCase{
		logger:      logger,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		objectRepo:  objectRepo,
		assetRepo:   assetRepo,
	}
}

func (uc *PostUserMeUpdateAvatarUseCase) Execute(
	ctx context.Context,
	input dto.PostUserMeUpdateAvatarInput,
) (dto.PostUserMeUpdateAvatarOutput, *dto.Error) {
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContextKeyRequestID)
	log.Debug("Use case called")

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/post_user_me_update_avatar",
			entity.ErrPostUserMeUpdateAvatarParamsInvalid,
			err.Error(),
		)
		return dto.PostUserMeUpdateAvatarOutput{}, &derr
	}

	// Get user ID from user session
	userID, err := uc.sessionRepo.GetUserIDByToken(ctx, input.AccessToken)
	if err != nil {
		derr := dto.NewError(
			"usecase/post_user_me_update_avatar",
			entity.ErrSessionNotFound,
			"failed to get user ID from access token",
		)
		log.Error("Failed to get user ID from access token",
			log.ToAny("error", err),
		)
		return dto.PostUserMeUpdateAvatarOutput{}, &derr
	}

	// Upload avatar file to object storage
	bucketName := "avatars"

	// The key is formed as "<userID>_<timestamp>.<format>"
	format := input.MimeFormat[len("image/"):] // govalidator ensures correct format
	key := fmt.Sprintf("%d/%d.%s", userID, time.Now().Unix(), format)

	// Upload avatar file to object storage
	avatarS3Key, err := uc.objectRepo.UploadObject(ctx, bucketName, key, input.MimeFormat, input.Bytes)
	if err != nil {
		derr := dto.NewError(
			"usecase/post_user_me_update_avatar",
			entity.ErrPostUserMeUpdateAvatarParamsInvalid,
			"failed to upload avatar file to object storage",
		)
		log.Error("Failed to upload avatar file to object storage",
			log.ToAny("error", err),
		)
		return dto.PostUserMeUpdateAvatarOutput{}, &derr
	}

	// Add avatar asset to asset repository
	asset := entity.Asset{
		S3Key:      avatarS3Key.GetPath(),
		MimeType:   input.MimeFormat,
		FileSizeMB: input.FileSizeMB,
	}

	// Create asset in asset repository
	assetID, err := uc.assetRepo.CreateAsset(ctx, asset)
	if err != nil {
		derr := dto.NewError(
			"usecase/post_user_me_update_avatar",
			entity.ErrPostUserMeUpdateAvatarParamsInvalid,
			"failed to create avatar asset in asset repository",
		)
		log.Error("Failed to create avatar asset in asset repository",
			log.ToAny("error", derr),
		)
		return dto.PostUserMeUpdateAvatarOutput{}, &derr
	}

	// TODO: extract image resolution (width, height) if needed
	assetImage := entity.AssetImage{
		AssetID:          assetID,
		ResolutionWidth:  100,
		ResolutionHeight: 100,
	}

	// Create asset image in asset repository
	assetImageID, err := uc.assetRepo.CreateAssetImage(ctx, assetImage)
	if err != nil {
		derr := dto.NewError(
			"usecase/post_user_me_update_avatar",
			entity.ErrPostUserMeUpdateAvatarParamsInvalid,
			"failed to create avatar asset image in asset repository",
		)
		log.Error("Failed to create avatar asset image in asset repository",
			log.ToAny("error", derr),
		)
		return dto.PostUserMeUpdateAvatarOutput{}, &derr
	}

	// Update user avatar s3 key in database
	err = uc.userRepo.UpdateUserAvatarKey(ctx, userID, assetImageID)
	if err != nil {
		derr := dto.NewError(
			"usecase/post_user_me_update_avatar",
			entity.ErrPostUserMeUpdateAvatarParamsInvalid,
			"failed to update user avatar s3 key in database",
		)
		log.Error("Failed to update user avatar s3 key in database",
			log.ToAny("error", derr),
		)
		return dto.PostUserMeUpdateAvatarOutput{}, &derr
	}

	// Get user avatar URL
	avatarURL, err := uc.objectRepo.GeneratePublicURL(ctx, bucketName, key)
	if err != nil {
		derr := dto.NewError(
			"usecase/post_user_me_update_avatar",
			entity.ErrPostUserMeUpdateAvatarParamsInvalid,
			"failed to generate user avatar URL",
		)
		log.Error("Failed to generate user avatar URL",
			log.ToAny("error", derr),
		)
		return dto.PostUserMeUpdateAvatarOutput{}, &derr
	}

	log.Debug("Use case completed successfully")
	return dto.PostUserMeUpdateAvatarOutput{
		URL: avatarURL.URL,
	}, nil

}
