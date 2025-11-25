package http

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type GetGenreUseCase struct {
	log             logger.Logger
	genreRepo       GenreRepository
	mediaRepo       MediaRepository
	getMediaUsecase *GetMediaUseCase
}

func NewGetGenreUseCase(
	log logger.Logger,
	genreRepo GenreRepository,
	mediaRepo MediaRepository,
	getMediaUsecase *GetMediaUseCase,
) *GetGenreUseCase {
	if log == nil {
		panic("log is nil")
	}
	if genreRepo == nil {
		panic("genreRepo is nil")
	}
	if mediaRepo == nil {
		panic("mediaRepo is nil")
	}
	if getMediaUsecase == nil {
		panic("getMediaUsecase is nil")
	}
	return &GetGenreUseCase{
		log:             log,
		genreRepo:       genreRepo,
		mediaRepo:       mediaRepo,
		getMediaUsecase: getMediaUsecase,
	}
}

func (uc *GetGenreUseCase) Execute(
	ctx context.Context,
	input dto.GetGenreInput,
) (
	dto.GetGenreOutput,
	*dto.Error,
) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.log, ctx, common.ContextKeyRequestID)
	log.Debug("GetGenreUseCase called",
		log.ToInt("genre_id", int(input.GenreID)),
	)

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_genre",
			entity.ErrGetGenreInvalidParams,
			err.Error(),
		)
		return dto.GetGenreOutput{}, &derr
	}

	// Get genre by ID
	genre, err := uc.genreRepo.GetGenreByID(ctx, input.GenreID)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_genre",
			entity.ErrGetGenreRepo,
			err.Error(),
		)
		return dto.GetGenreOutput{}, &derr
	}
	if genre == nil {
		derr := dto.NewError(
			"usecase/get_genre",
			entity.ErrGetGenreNotFound,
			"genre not found",
		)
		return dto.GetGenreOutput{}, &derr
	}

	// Get media IDs related to specific genre ID
	mediaIDs, err := uc.mediaRepo.GetMediasByGenreID(ctx, input.MediaLimit, input.MediaOffset, input.GenreID)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_genre",
			entity.ErrGetGenreRepo,
			err.Error(),
		)
		return dto.GetGenreOutput{}, &derr
	}

	// Get media entities by IDs
	medias := make([]dto.GetMediaOutput, 0, len(mediaIDs))
	for _, mediaID := range mediaIDs {
		getMediaOutput, derr := uc.getMediaUsecase.Execute(ctx, dto.GetMediaInput{
			MediaID: mediaID,
		})
		if derr != nil {
			log.Error("GetGenreUseCase failed to get media by ID", log.ToError(err), log.ToInt("media_id", int(mediaID)))
			continue
		}
		medias = append(medias, getMediaOutput)
	}

	output := dto.GetGenreOutput{
		Genre:  *genre,
		Medias: medias,
	}

	return output, nil
}
