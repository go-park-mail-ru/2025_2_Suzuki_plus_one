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
}

func NewGetGenreUseCase(
	log logger.Logger,
	genreRepo GenreRepository,
	mediaRepo MediaRepository,
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

	return &GetGenreUseCase{
		log:             log,
		genreRepo:       genreRepo,
		mediaRepo:       mediaRepo,
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

	output := dto.GetGenreOutput{
		Genre:  *genre,
	}

	return output, nil
}
