package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type GetActorMediaUseCase struct {
	logger          logger.Logger
	actorRepo       ActorRepository
	getMediaUseCase *GetMediaUseCase
}

func NewGetActorMediaUseCase(
	logger logger.Logger,
	actorRepo ActorRepository,
	getMediaUseCase *GetMediaUseCase,
) *GetActorMediaUseCase {
	return &GetActorMediaUseCase{
		logger:          logger,
		actorRepo:       actorRepo,
		getMediaUseCase: getMediaUseCase,
	}
}

func (uc *GetActorMediaUseCase) Execute(ctx context.Context, input dto.GetActorMediaInput) (dto.GetActorMediaOutput, *dto.Error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContextKeyRequestID)

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_actor_media",
			err,
			"Invalid get actor input parameters",
		)
		return dto.GetActorMediaOutput{}, &derr
	}

	// Get medias associated with the actor
	medias, err := uc.actorRepo.GetMediasByActorID(ctx, input.ActorID)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_actor_media",
			err,
			"Failed to get medias by actor ID",
		)
		log.Error("Failed to get medias by actor ID", log.ToError(err))
		return dto.GetActorMediaOutput{}, &derr
	}

	// Map actor and medias to output DTO
	output := dto.GetActorMediaOutput{
		Medias: make([]dto.GetMediaOutput, 0, len(medias)),
	}

	// Get medias for output actor
	for _, media := range medias {
		mediaOutput, err := uc.getMediaUseCase.Execute(ctx, dto.GetMediaInput{MediaID: media.MediaID})
		if err != nil {
			log.Error("Failed to get media output", err)
			continue
		}
		output.Medias = append(output.Medias, mediaOutput)
	}
	return output, nil
}
