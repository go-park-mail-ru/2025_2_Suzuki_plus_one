package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type GetActorUseCase struct {
	logger           logger.Logger
	actorRepo        ActorRepository
	getMediaUseCase  *GetMediaUseCase
	getObjectUseCase *GetObjectUseCase
}

func NewGetActorUseCase(
	logger logger.Logger,
	actorRepo ActorRepository,
	getMediaUseCase *GetMediaUseCase,
	getObjectUseCase *GetObjectUseCase,
) *GetActorUseCase {
	return &GetActorUseCase{
		logger:           logger,
		actorRepo:        actorRepo,
		getMediaUseCase:  getMediaUseCase,
		getObjectUseCase: getObjectUseCase,
	}
}

func (uc *GetActorUseCase) Execute(ctx context.Context, input dto.GetActorInput) (dto.GetActorOutput, *dto.Error) {
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
		uc.logger.Error("Failed to get actor by ID", uc.logger.ToError(err))
		return dto.GetActorOutput{}, &derr
	}

	// Get medias associated with the actor
	medias, err := uc.actorRepo.GetMediasByActorID(ctx, input.ActorID)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_actor",
			err,
			"Failed to get medias by actor ID",
		)
		uc.logger.Error("Failed to get medias by actor ID", uc.logger.ToError(err))
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
		uc.logger.Error("Failed to get actor images by actor ID", uc.logger.ToError(err))
		return dto.GetActorOutput{}, &derr
	}

	// Map actor and medias to output DTO
	output := dto.GetActorOutput{
		Actor:     *actor,
		Medias:    make([]dto.GetMediaOutput, 0, len(medias)),
		ImageURLs: make([]string, 0, len(imageS3Keys)),
	}

	// Get medias for output actor
	for _, media := range medias {
		mediaOutput, err := uc.getMediaUseCase.Execute(ctx, dto.GetMediaInput{MediaID: media.MediaID})
		if err != nil {
			uc.logger.Error("Failed to get media output", err)
			continue
		}
		output.Medias = append(output.Medias, mediaOutput)
	}
	// Get image URLs for output actor
	for _, imageS3Key := range imageS3Keys {
		imageURL, err := uc.getObjectUseCase.Execute(ctx,
			dto.GetObjectInput{Key: imageS3Key.Key, BucketName: imageS3Key.BucketName})
		if err != nil {
			uc.logger.Error("Failed to get image URL", err)
			continue
		}
		output.ImageURLs = append(output.ImageURLs, imageURL.URL)
	}
	return output, nil
}
