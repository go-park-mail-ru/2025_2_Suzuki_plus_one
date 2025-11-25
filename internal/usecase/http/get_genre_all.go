package http

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

const GET_GENRE_ALL_LIMIT_MEDIA = 10
const GET_GENRE_ALL_OFFSET_MEDIA = 0

type GetGenreAllUseCase struct {
	log             logger.Logger
	genreRepo       GenreRepository
	getGenreUsecase *GetGenreUseCase
}

func NewGetGenreAllUseCase(
	log logger.Logger,
	genreRepo GenreRepository,
	getGenreUsecase *GetGenreUseCase,
) *GetGenreAllUseCase {
	if log == nil {
		panic("log is nil")
	}
	if genreRepo == nil {
		panic("genreRepo is nil")
	}
	if getGenreUsecase == nil {
		panic("getGenreUsecase is nil")
	}
	return &GetGenreAllUseCase{
		log:             log,
		genreRepo:       genreRepo,
		getGenreUsecase: getGenreUsecase,
	}
}

func (uc *GetGenreAllUseCase) Execute(
	ctx context.Context,
	input dto.GetGenreAllInput,
) (
	dto.GetGenreAllOutput,
	*dto.Error,
) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.log, ctx, common.ContextKeyRequestID)
	log.Debug("GetGenreAllUseCase called")

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_genre_all",
			entity.ErrGetGenreAllInvalidParams,
			err.Error(),
		)
		return dto.GetGenreAllOutput{}, &derr
	}

	// Get all genre IDs
	genreIDs, err := uc.genreRepo.GetAllGenreIDs(ctx)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_genre_all",
			entity.ErrGetAllGenresFailed,
			err.Error(),
		)
		return dto.GetGenreAllOutput{}, &derr
	}

	// Execute GetGenreUseCase for each genre ID
	var genres []dto.GetGenreOutput
	for _, genreID := range genreIDs {
		genreOutput, derr := uc.getGenreUsecase.Execute(ctx, dto.GetGenreInput{
			GenreID: genreID,
			MediaLimit: GET_GENRE_ALL_LIMIT_MEDIA,
			MediaOffset: GET_GENRE_ALL_OFFSET_MEDIA,
		})
		if derr != nil {
			return dto.GetGenreAllOutput{}, derr
		}
		genres = append(genres, genreOutput)
	}

	output := dto.GetGenreAllOutput{
		Genres: genres,
	}

	return output, nil
}
