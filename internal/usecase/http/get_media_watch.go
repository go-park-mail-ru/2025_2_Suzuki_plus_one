package http

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type GetMediaWatchUseCase struct {
	logger           logger.Logger
	mediaRepo        MediaRepository
	getObjectUseCase *GetObjectUseCase
	userRepo 	  UserRepository
}

func NewGetMediaWatchUseCase(
	logger logger.Logger,
	mediaRepo MediaRepository,
	getObjectUseCase *GetObjectUseCase,
	userRepo UserRepository,
) *GetMediaWatchUseCase {
	if logger == nil {
		panic("NewGetMediaWatchUseCase: logger is nil")
	}
	if mediaRepo == nil {
		panic("NewGetMediaWatchUseCase: mediaRepo is nil")
	}
	if getObjectUseCase == nil {
		panic("NewGetMediaWatchUseCase: getObjectUseCase is nil")
	}
	if userRepo == nil {
		panic("NewGetMediaWatchUseCase: userRepo is nil")
	}
	return &GetMediaWatchUseCase{
		logger:           logger,
		mediaRepo:        mediaRepo,
		getObjectUseCase: getObjectUseCase,
		userRepo:         userRepo,
	}
}

func (uc *GetMediaWatchUseCase) Execute(ctx context.Context, input dto.GetMediaWatchInput) (dto.GetMediaWatchOutput, *dto.Error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContextKeyRequestID)

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_media_watch",
			entity.ErrGetMediaParamsInvalid,
			err.Error(),
		)
		return dto.GetMediaWatchOutput{}, &derr
	}

	// Get media watch URL from database
	s3Key, err := uc.mediaRepo.GetMediaWatchKey(ctx, input.MediaID)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_media_watch",
			err,
			"Failed to get media watch URL by ID",
		)
		log.Error("Failed to get media watch URL by ID", log.ToError(err))
		return dto.GetMediaWatchOutput{}, &derr
	}


	media, err := uc.mediaRepo.GetMediaByID(ctx, input.MediaID)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_media_watch",
			err,
			"Failed to get media info by ID",
		)
		log.Error("Failed to get media info by ID", log.ToError(err))
		return dto.GetMediaWatchOutput{}, &derr
	}

	// If media is movie or series, validate access token and user subscription status
	if media.MediaType == "episode" {
		// Get user ID from access token
		userID, err := common.ValidateToken(input.AccessToken)
		if err != nil {
			derr := dto.NewError(
				"usecase/get_media_watch",
				entity.ErrGetMediaAccessDenied,
				"access token is invalid: "+err.Error(),
			)
			log.Error("Access token is invalid", log.ToError(err))
			return dto.GetMediaWatchOutput{}, &derr
		}

		// Get user subscription status
		status, err := uc.userRepo.GetUserSubscriptionStatus(ctx, userID)
		if err != nil {
			derr := dto.NewError(
				"usecase/get_media_watch",
				err,
				"Failed to get user subscription status",
			)
			log.Error("Failed to get user subscription status", log.ToError(err))
			return dto.GetMediaWatchOutput{}, &derr
		}

		// Check if subscription status is active
		if status != "active" {
			derr := dto.NewError(
				"usecase/get_media_watch",
				entity.ErrGetMediaAccessDenied,
				"user subscription is not active",
			)
			log.Error("User subscription is not active")
			return dto.GetMediaWatchOutput{}, &derr
		}
	}
	
	// Get presigned URL from object use case
	object, derr := uc.getObjectUseCase.Execute(ctx, dto.GetObjectInput{
		Key:        s3Key.Key,
		BucketName: s3Key.BucketName,
	})
	if derr != nil {
		log.Error("Failed to get presigned URL for media watch", derr)
		return dto.GetMediaWatchOutput{}, derr
	}

	return dto.GetMediaWatchOutput{
		URL: object.URL,
	}, nil
}
