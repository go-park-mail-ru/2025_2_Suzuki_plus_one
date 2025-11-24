package usecase

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type GetObjectUseCase struct {
	logger     logger.Logger
	objectRepo ObjectRepository
}

func NewGetObjectUseCase(
	logger logger.Logger,
	objectRepo ObjectRepository,
) *GetObjectUseCase {
	return &GetObjectUseCase{
		logger:     logger,
		objectRepo: objectRepo,
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
		log.Error("Invalid get object input parameters", log.ToError(err))
		return dto.GetObjectOutput{}, &derr
	}
	log.Debug("Get object input parameters are valid", log.ToAny("input", input))
	// Get object from repository

	// Differentiate between public and private buckets
	// For now only "medias" bucket is private
	var object *entity.URL
	var err error
	if input.BucketName == "medias" {
		// TODO: linkAliveDuration is hardcoded here, can be changed later if needed
		linkAliveDuration := time.Minute * 15
		object, err = uc.objectRepo.GeneratePresignedURL(ctx, input.BucketName, input.Key, linkAliveDuration)
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
