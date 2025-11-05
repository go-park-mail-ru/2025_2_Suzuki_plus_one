package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type GetMediaWatchUseCase struct {
	logger           logger.Logger
	mediaRepo        MediaRepository
	getObjectUseCase *GetObjectUseCase
}

func NewGetMediaWatchUseCase(
	logger logger.Logger,
	mediaRepo MediaRepository,
	getObjectUseCase *GetObjectUseCase,
) *GetMediaWatchUseCase {
	return &GetMediaWatchUseCase{
		logger:           logger,
		mediaRepo:        mediaRepo,
		getObjectUseCase: getObjectUseCase,
	}
}

func (uc *GetMediaWatchUseCase) Execute(ctx context.Context, input dto.GetMediaWatchInput) (dto.GetMediaWatchOutput, *dto.Error) {
	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_media",
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
		uc.logger.Error("Failed to get media watch URL by ID", uc.logger.ToError(err))
		return dto.GetMediaWatchOutput{}, &derr
	}

	// Get presigned URL from object use case
	object, derr := uc.getObjectUseCase.Execute(ctx, dto.GetObjectInput{
		Key:        s3Key.Key,
		BucketName: s3Key.BucketName,
	})
	if derr != nil {
		uc.logger.Error("Failed to get presigned URL for media watch", derr)
		return dto.GetMediaWatchOutput{}, derr
	}

	return dto.GetMediaWatchOutput{
		URL: object.URL,
	}, nil
}
