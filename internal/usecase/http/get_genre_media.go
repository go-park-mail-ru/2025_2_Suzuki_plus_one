package http

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type GetGenreMediaUseCase struct {
	logger          logger.Logger
	mediaRepo       MediaRepository
	getMediaUseCase *GetMediaUseCase
}

func NewGetGenreMediaUseCase(
	logger logger.Logger,
	mediaRepo MediaRepository,
	getMediaUseCase *GetMediaUseCase,
) *GetGenreMediaUseCase {
	if logger == nil {
		panic("logger is nil")
	}
	if mediaRepo == nil {
		panic("mediaRepo is nil")
	}
	if getMediaUseCase == nil {
		panic("getMediaUseCase is nil")
	}
	return &GetGenreMediaUseCase{
		logger:          logger,
		mediaRepo:       mediaRepo,
		getMediaUseCase: getMediaUseCase,
	}
}

func (uc *GetGenreMediaUseCase) Execute(ctx context.Context, input dto.GetGenreMediaInput) (dto.GetGenreMediaOutput, *dto.Error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContextKeyRequestID)

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_genre_media",
			err,
			"Invalid get genre media input parameters",
		)
		return dto.GetGenreMediaOutput{}, &derr
	}

	// Get media IDs related to specific genre ID
	mediaIDs, err := uc.mediaRepo.GetMediasByGenreID(ctx, input.Limit, input.Offset, input.GenreID)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_genre_media",
			entity.ErrGetGenreRepo,
			err.Error(),
		)
		return dto.GetGenreMediaOutput{}, &derr
	}

	// Get media entities by IDs
	medias := make([]dto.GetMediaOutput, 0, len(mediaIDs))
	for _, mediaID := range mediaIDs {
		getMediaOutput, derr := uc.getMediaUseCase.Execute(ctx, dto.GetMediaInput{
			MediaID: mediaID,
		})
		if derr != nil {
			log.Error("GetGenreUseCase failed to get media by ID", log.ToError(err), log.ToInt("media_id", int(mediaID)))
			continue
		}
		medias = append(medias, getMediaOutput)
	}

	output := dto.GetGenreMediaOutput{
		Medias: medias,
	}

	return output, nil
}
