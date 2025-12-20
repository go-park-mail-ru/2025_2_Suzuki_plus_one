package http

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type GetObjectUseCase struct {
	logger              logger.Logger
	objectRepo          ObjectRepository
	presignedExpiration time.Duration
}

func NewGetObjectUseCase(
	logger logger.Logger,
	objectRepo ObjectRepository,
	presignedExpiration time.Duration,
) *GetObjectUseCase {
	if logger == nil {
		panic("logger is nil")
	}
	if objectRepo == nil {
		panic("objectRepo is nil")
	}
	if presignedExpiration <= 0 {
		panic("presignedExpiration is invalid")
	}
	return &GetObjectUseCase{
		logger:              logger,
		objectRepo:          objectRepo,
		presignedExpiration: presignedExpiration,
	}
}

func (uc *GetObjectUseCase) Execute(ctx context.Context, input dto.GetObjectInput) (dto.GetObjectOutput, *dto.Error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContextKeyRequestID)

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_object",
			entity.ErrGetObjectParamsInvalid,
			err.Error(),
		)
		log.Error("Invalid get object input parameters", log.ToError(err), log.ToAny("input", input))
		return dto.GetObjectOutput{}, &derr
	}
	log.Debug("Get object input parameters are valid", log.ToAny("input", input))
	// Get object from repository

	// Differentiate between public and private buckets
	// For now only "medias" bucket is private
	var object *entity.URL
	var err error
	if input.BucketName == "medias" {
		object, err = uc.objectRepo.GeneratePresignedURL(ctx, input.BucketName, input.Key, uc.presignedExpiration)
	} else {
		object, err = uc.objectRepo.GeneratePublicURL(ctx, input.BucketName, input.Key)
	}

	if err != nil {
		derr := dto.NewError(
			"usecase/get_object",
			entity.ErrGetObjectFailed,
			err.Error(),
		)
		log.Error("Failed to get object", log.ToError(err))
		return dto.GetObjectOutput{}, &derr
	}

	// Return output DTO
	return dto.GetObjectOutput{
		URL: object.URL,
	}, nil
}
