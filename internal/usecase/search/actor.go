package search

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type SearchActorUsecase struct {
	logger    logger.Logger
	actorRepo ActorRepository
}

func NewSearchActorUsecase(logger logger.Logger, actorRepo ActorRepository) *SearchActorUsecase {
	if actorRepo == nil {
		panic("actorRepo is nil")
	}
	return &SearchActorUsecase{
		logger:    logger,
		actorRepo: actorRepo,
	}
}

func (u *SearchActorUsecase) Execute(ctx context.Context, input dto.GetSearchInput) (dto.GetSearchOutput, *dto.Error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(u.logger, ctx, common.ContextKeyRequestID)
	log.Debug("Executing SearchActorUsecase with input: ", log.ToString("qiery", input.Query))

	//
	actorIDs, err := u.actorRepo.SearchActor(ctx, input.Query, input.Limit, input.Offset)
	if err != nil {
		log.Error("SearchActorUsecase: failed to search actors: ", log.ToError(err))
		derr := dto.NewError("search/get_search", entity.ErrorSearchActorFailed, "failed to search actors")
		return dto.GetSearchOutput{}, &derr
	}

	// Prepare output
	output := dto.GetSearchOutput{
		Medias: make([]dto.GetMediaOutput, 0, len(actorIDs)),
	}

	// Retrieve media Ids, skip the rest of the details
	for _, actorID := range actorIDs {
		actor := dto.GetActorOutput{}
		actor.ID = actorID
		output.Actors = append(output.Actors, actor)
	}

	return output, nil

}
