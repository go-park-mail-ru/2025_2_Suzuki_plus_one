package usecase

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type GetObjectUsecase struct {
	logger     logger.Logger
	objectRepo ObjectRepository
}

func NewGetObjectUsecase(
	logger logger.Logger,
	objectRepo ObjectRepository,
) *GetObjectUsecase {
	return &GetObjectUsecase{
		logger:     logger,
		objectRepo: objectRepo,
	}
}

func (uc *GetObjectUsecase) Execute(ctx context.Context, input dto.GetObjectInput) (dto.GetObjectOutput, *dto.Error) {
	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_object",
			entity.ErrGetObjectParamsInvalid,
			err.Error(),
		)
		return dto.GetObjectOutput{}, &derr
	}

	// Get object from repository
	// TODO: linkAliveDuration is hardcoded here, can be changed later if needed
	linkAliveDuration := time.Minute * 15
	object, err := uc.objectRepo.GetObject(ctx, input.Key, input.BucketName, linkAliveDuration)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_object",
			entity.ErrGetObjectFailed,
			err.Error(),
		)
		return dto.GetObjectOutput{}, &derr
	}

	// Return output DTO
	return dto.GetObjectOutput{
		URL: object.URL,
	}, nil
}
