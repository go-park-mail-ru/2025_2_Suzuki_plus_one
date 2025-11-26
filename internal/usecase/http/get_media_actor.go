package http

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type GetMediaActorUseCase struct {
	logger          logger.Logger
	actorRepo       ActorRepository
	getActorUseCase *GetActorUseCase
}

func NewGetMediaActorUseCase(
	logger logger.Logger,
	actorRepo ActorRepository,
	getActorUseCase *GetActorUseCase,
) *GetMediaActorUseCase {
	return &GetMediaActorUseCase{
		logger:          logger,
		actorRepo:       actorRepo,
		getActorUseCase: getActorUseCase,
	}
}

func (uc *GetMediaActorUseCase) Execute(ctx context.Context, input dto.GetMediaActorInput) (dto.GetMediaActorOutput, *dto.Error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContextKeyRequestID)

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_media_actor",
			entity.ErrGetMediaActorParamsInvalid,
			err.Error(),
		)
		return dto.GetMediaActorOutput{}, &derr
	}

	// Get actors associated with the media
	actors, err := uc.actorRepo.GetActorsByMediaID(ctx, input.MediaID)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_media_actor",
			err,
			"Failed to get actors by media ID",
		)
		log.Error("Failed to get actors by media ID", log.ToError(err))
		return dto.GetMediaActorOutput{}, &derr
	}

	// Get actors
	actorsDTO := make([]dto.GetActorOutput, 0, len(actors))
	for _, actor := range actors {
		actorOutput, err := uc.getActorUseCase.Execute(ctx, dto.GetActorInput{ActorID: actor.ID})
		if err != nil {
			log.Error("Failed to get actor output", log.ToAny("error", err))
			continue
		}
		actorsDTO = append(actorsDTO, actorOutput)
	}

	return dto.GetMediaActorOutput{Actors: actorsDTO}, nil
}
