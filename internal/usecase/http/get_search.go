package http

import (
	"context"
	"strings"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"

	unidecode "github.com/mozillazg/go-unidecode"
)

type GetSearchUseCase struct {
	logger          logger.Logger
	searchService   ServiceSearchRepository
	getMediaUsecase *GetMediaUseCase
	getActorUsecase *GetActorUseCase
}

func NewGetSearchUseCase(
	logger logger.Logger,
	searchService ServiceSearchRepository,
	getMediaUsecase *GetMediaUseCase,
	getActorUsecase *GetActorUseCase,
) *GetSearchUseCase {
	if logger == nil {
		panic("nil logger passed to NewGetSearchUseCase")
	}
	if searchService == nil {
		panic("nil searchService passed to NewGetSearchUseCase")
	}
	if getMediaUsecase == nil {
		panic("nil getMediaUsecase passed to NewGetSearchUseCase")
	}
	if getActorUsecase == nil {
		panic("nil getActorUsecase passed to NewGetSearchUseCase")
	}
	return &GetSearchUseCase{
		logger:          logger,
		searchService:   searchService,
		getMediaUsecase: getMediaUsecase,
		getActorUsecase: getActorUsecase,
	}
}

func (uc *GetSearchUseCase) Execute(ctx context.Context, input dto.GetSearchInput) (dto.GetSearchOutput, *dto.Error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContextKeyRequestID)

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_object",
			entity.ErrGetSearchParamsInvalid,
			err.Error(),
		)
		log.Error("Invalid get object input parameters", log.ToError(err))
		return dto.GetSearchOutput{}, &derr
	}
	log.Debug("Get object input parameters are valid", log.ToAny("input", input))

	// Optional: collapse multiple spaces
	query := strings.Join(strings.Fields(input.Query), " ")

	// Transliterate query
	query = unidecode.Unidecode(query)
	log.Debug("Transliterated search query", log.ToString("original", input.Query), log.ToString("transliterated", query))

	// Call Search service to get IDs
	var media_ids []uint
	var actor_ids []uint
	var err error
	switch input.Type {
	case "media":
		media_ids, err = uc.searchService.CallSearchMedia(ctx, query, input.Limit, input.Offset)
	case "actor":
		actor_ids, err = uc.searchService.CallSearchActors(ctx, query, input.Limit, input.Offset)
	case "any":
		var err1, err2 error
		media_ids, err1 = uc.searchService.CallSearchMedia(ctx, query, input.Limit, input.Offset)
		actor_ids, err2 = uc.searchService.CallSearchActors(ctx, query, input.Limit, input.Offset)
		if err1 != nil {
			err = err1
		} else if err2 != nil {
			err = err2
		}
	default:
		derr := dto.NewError(
			"usecase/get_search",
			entity.ErrGetSearchParamsInvalid,
			"invalid search type: "+input.Type,
		)
		log.Error("Invalid search type parameter", log.ToString("type", input.Type))
		return dto.GetSearchOutput{}, &derr
	}
	if err != nil {
		derr := dto.NewError(
			"usecase/get_search",
			entity.ErrGetSearchSearchServiceFailed,
			"search service failed: "+err.Error(),
		)
		log.Error("Search service call failed", log.ToError(err))
		return dto.GetSearchOutput{}, &derr
	}

	// Get full media details
	medias := make([]dto.GetMediaOutput, 0, len(media_ids))
	for _, media_id := range media_ids {
		mediaOutput, derr := uc.getMediaUsecase.Execute(ctx, dto.GetMediaInput{MediaID: media_id})
		if derr != nil {
			log.Error("Failed to get media details", log.ToAny("media_id", media_id), log.ToString("error", derr.Message))
			continue
		}
		medias = append(medias, mediaOutput)
	}

	// Get full actor details
	actors := make([]dto.GetActorOutput, 0, len(actor_ids))
	for _, actor_id := range actor_ids {
		actorOutput, derr := uc.getActorUsecase.Execute(ctx, dto.GetActorInput{ActorID: actor_id})
		if derr != nil {
			log.Error("Failed to get actor details", log.ToAny("actor_id", actor_id), log.ToString("error", derr.Message))
			continue
		}
		actors = append(actors, actorOutput)
	}

	// Return output DTO
	return dto.GetSearchOutput{
		Medias: medias,
		Actors: actors,
	}, nil
}
