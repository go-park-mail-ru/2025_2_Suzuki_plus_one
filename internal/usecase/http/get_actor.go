package http

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type GetActorUseCase struct {
	logger           logger.Logger
	actorRepo        ActorRepository
	getObjectUseCase *GetObjectUseCase
}

func NewGetActorUseCase(
	logger logger.Logger,
	actorRepo ActorRepository,
	getObjectUseCase *GetObjectUseCase,
) *GetActorUseCase {
	return &GetActorUseCase{
		logger:           logger,
		actorRepo:        actorRepo,
		getObjectUseCase: getObjectUseCase,
	}
}

func (uc *GetActorUseCase) Execute(ctx context.Context, input dto.GetActorInput) (dto.GetActorOutput, *dto.Error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContextKeyRequestID)

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_actor",
			err,
			"Invalid get actor input parameters",
		)
		return dto.GetActorOutput{}, &derr
	}

	// Get actor from repository
	actor, err := uc.actorRepo.GetActorByID(ctx, input.ActorID)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_actor",
			err,
			"Failed to get actor by ID",
		)
		log.Error("Failed to get actor by ID", log.ToError(err))
		return dto.GetActorOutput{}, &derr
	}

	// Get actor image URLs
	imageS3Keys, err := uc.actorRepo.GetActorImageS3(ctx, input.ActorID)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_actor",
			err,
			"Failed to get actor images by actor ID",
		)
		log.Error("Failed to get actor images by actor ID", log.ToError(err))
		return dto.GetActorOutput{}, &derr
	}

	// Map actor and medias to output DTO
	output := dto.GetActorOutput{
		Actor:     *actor,
		ImageURLs: make([]string, 0, len(imageS3Keys)),
	}

	// Get image URLs for output actor
	for _, imageS3Key := range imageS3Keys {
		imageURL, err := uc.getObjectUseCase.Execute(ctx,
			dto.GetObjectInput{Key: imageS3Key.Key, BucketName: imageS3Key.BucketName})
		if err != nil {
			log.Error("Failed to get image URL", err)
			continue
		}
		output.ImageURLs = append(output.ImageURLs, imageURL.URL)
	}
	return output, nil
}
