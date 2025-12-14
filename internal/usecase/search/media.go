package search

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type SearchMediaUsecase struct {
	logger    logger.Logger
	mediaRepo MediaRepository
}

func NewSearchMediaUsecase(logger logger.Logger, mediaRepo MediaRepository) *SearchMediaUsecase {
	if mediaRepo == nil {
		panic("mediaRepo is nil")
	}
	return &SearchMediaUsecase{
		logger:    logger,
		mediaRepo: mediaRepo,
	}
}

func (u *SearchMediaUsecase) Execute(ctx context.Context, input dto.GetSearchInput) (dto.GetSearchOutput, *dto.Error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(u.logger, ctx, common.ContextKeyRequestID)
	log.Debug("Executing SearchMediaUsecase with input: ", log.ToString("qiery", input.Query))

	//
	mediaIDs, err := u.mediaRepo.SearchMedia(ctx, input.Query, input.Limit, input.Offset)
	if err != nil {
		log.Error("SearchMediaUsecase: failed to search media: ", log.ToError(err))
		derr := dto.NewError("search/get_search", entity.ErrorSearchMediaFailed, "failed to search media")
		return dto.GetSearchOutput{}, &derr
	}

	// Prepare output
	output := dto.GetSearchOutput{
		Medias: make([]dto.GetMediaOutput, 0, len(mediaIDs)),
	}

	// Retrieve media Ids, skip the rest of the details
	for _, mediaID := range mediaIDs {
		media := dto.GetMediaOutput{}
		media.MediaID = mediaID
		output.Medias = append(output.Medias, media)
	}

	return output, nil

}
