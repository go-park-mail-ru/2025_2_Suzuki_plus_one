package usecase

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

const GET_MOVIE_RECOMMENDATION_LIMIT_MIN = 1
const GET_MOVIE_RECOMMENDATION_LIMIT_MAX = 20

type GetMovieRecommendationsUsecase struct {
	logger          logger.Logger
	movieRepo       MediaRepository
	getMediaUseCase *GetMediaUseCase
}

func NewGetMovieRecommendationsUsecase(
	logger logger.Logger,
	movieRepo MediaRepository,
	getMediaUseCase *GetMediaUseCase,
) *GetMovieRecommendationsUsecase {
	return &GetMovieRecommendationsUsecase{
		logger:          logger,
		movieRepo:       movieRepo,
		getMediaUseCase: getMediaUseCase,
	}
}

func (uc *GetMovieRecommendationsUsecase) Execute(
	ctx context.Context,
	input dto.GetMovieRecommendationsInput,
) (
	dto.GetMovieRecommendationsOutput,
	*dto.Error,
) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContexKeyRequestID)

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_movie_recommendations",
			entity.ErrGetMovieRecommendationsParamsInvalid,
			err.Error(),
		)
		return dto.GetMovieRecommendationsOutput{}, &derr
	}

	// Get total movie number from repository
	movie_number, err := uc.movieRepo.GetMediaCount(ctx, "movie")
	if err != nil {
		derr := dto.NewError("adapter/get_movie_recommendations", err, "Can't get movie count")
		log.Error("Can't get movie count", log.ToError(err))
		return dto.GetMovieRecommendationsOutput{}, &derr
	}

	// Validate ?limit boundaries
	if input.Limit < GET_MOVIE_RECOMMENDATION_LIMIT_MIN || input.Limit > GET_MOVIE_RECOMMENDATION_LIMIT_MAX {
		derr := dto.NewError(
			"adapter/get_movie_recommendations",
			entity.ErrGetMovieRecommendationsCantReturnAll,
			fmt.Sprintf(
				"limit must be between %d and %d (not %d)",
				GET_MOVIE_RECOMMENDATION_LIMIT_MIN,
				GET_MOVIE_RECOMMENDATION_LIMIT_MAX,
				input.Limit,
			),
		)
		log.Error("Invalid limit for movie recommendations", log.ToError(fmt.Errorf("limit: %d", input.Limit)))
		return dto.GetMovieRecommendationsOutput{}, &derr
	}

	// Return error if limit+offset exceeds total number of movies
	if input.Offset+input.Limit > uint(movie_number) {
		derr := dto.NewError(
			"usecase/get_movie_recommendations",
			entity.ErrGetMovieRecommendationsParamsInvalid,
			fmt.Sprintf(
				"limit+offset out of movies range: %d+%d of %d",
				input.Offset,
				input.Limit,
				movie_number,
			),
		)
		log.Error("limit+offset out of movies range", log.ToError(fmt.Errorf("offset: %d, limit: %d, total: %d", input.Offset, input.Limit, movie_number)))
		return dto.GetMovieRecommendationsOutput{}, &derr
	}

	// Calculate recommendation movie IDs (dummy logic for now)
	movieIDs, err := uc.movieRepo.GetMediaRandomIds(ctx, input.Limit, input.Offset, "movie")
	if err != nil {
		derr := dto.NewError("adapter/get_movie_recommendations", err, "Can't get movie recommendations")
		log.Error("Can't get movie recommendations", log.ToError(err))
		return dto.GetMovieRecommendationsOutput{}, &derr
	}

	// Convert []entity.Movie to []dto.Movie
	dtoMovies := make([]dto.GetMediaOutput, 0, len(movieIDs))
	for _, movieID := range movieIDs {
		// OPTIMIZATION: Batch fetch movies by IDs to reduce number of queries
		movieDTO, derr := uc.getMediaUseCase.Execute(ctx, dto.GetMediaInput{MediaID: movieID})
		if derr != nil {
			log.Error("Failed to get movie by ID in recommendations", derr)
			return dto.GetMovieRecommendationsOutput{}, derr
		}
		dtoMovies = append(dtoMovies, movieDTO)
	}

	return dto.GetMovieRecommendationsOutput{Movies: dtoMovies}, nil
}
