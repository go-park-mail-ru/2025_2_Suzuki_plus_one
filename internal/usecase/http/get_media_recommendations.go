package http

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

type GetMediaRecommendationsUsecase struct {
	logger          logger.Logger
	movieRepo       MediaRepository
	getMediaUseCase *GetMediaUseCase
}

func NewGetMediaRecommendationsUsecase(
	logger logger.Logger,
	movieRepo MediaRepository,
	getMediaUseCase *GetMediaUseCase,
) *GetMediaRecommendationsUsecase {
	return &GetMediaRecommendationsUsecase{
		logger:          logger,
		movieRepo:       movieRepo,
		getMediaUseCase: getMediaUseCase,
	}
}

func (uc *GetMediaRecommendationsUsecase) Execute(
	ctx context.Context,
	input dto.GetMediaRecommendationsInput,
) (
	dto.GetMediaRecommendationsOutput,
	*dto.Error,
) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContextKeyRequestID)

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_movie_recommendations",
			entity.ErrGetMediaRecommendationsParamsInvalid,
			err.Error(),
		)
		return dto.GetMediaRecommendationsOutput{}, &derr
	}

	// Get total movie number from repository
	movie_number, err := uc.movieRepo.GetMediaCount(ctx, input.Type)
	if err != nil {
		derr := dto.NewError("adapter/get_movie_recommendations", err, "Can't get movie count")
		log.Error("Can't get movie count", log.ToError(err))
		return dto.GetMediaRecommendationsOutput{}, &derr
	}

	// Validate ?limit boundaries
	if input.Limit < GET_MOVIE_RECOMMENDATION_LIMIT_MIN || input.Limit > GET_MOVIE_RECOMMENDATION_LIMIT_MAX {
		derr := dto.NewError(
			"adapter/get_movie_recommendations",
			entity.ErrGetMediaRecommendationsCantReturnAll,
			fmt.Sprintf(
				"limit must be between %d and %d (not %d)",
				GET_MOVIE_RECOMMENDATION_LIMIT_MIN,
				GET_MOVIE_RECOMMENDATION_LIMIT_MAX,
				input.Limit,
			),
		)
		log.Error("Invalid limit for movie recommendations", log.ToError(fmt.Errorf("limit: %d", input.Limit)))
		return dto.GetMediaRecommendationsOutput{}, &derr
	}

	// Return error if limit+offset exceeds total number of movies
	if input.Offset+input.Limit > uint(movie_number) {
		derr := dto.NewError(
			"usecase/get_movie_recommendations",
			entity.ErrGetMediaRecommendationsParamsInvalid,
			fmt.Sprintf(
				"limit+offset > media count: %d+%d > %d",
				input.Offset,
				input.Limit,
				movie_number,
			),
		)
		log.Error("limit+offset out of movies range", log.ToError(fmt.Errorf("offset: %d, limit: %d, total: %d", input.Offset, input.Limit, movie_number)))
		return dto.GetMediaRecommendationsOutput{}, &derr
	}

	// Calculate recommendation movie IDs (dummy logic for now)
	movieIDs, err := uc.movieRepo.GetMediaSortedByName(ctx, input.Limit, input.Offset, input.Type)
	if err != nil {
		derr := dto.NewError("adapter/get_movie_recommendations", err, "Can't get movie recommendations")
		log.Error("Can't get movie recommendations", log.ToError(err))
		return dto.GetMediaRecommendationsOutput{}, &derr
	}

	// Convert []entity.Movie to []dto.Movie
	dtoMovies := make([]dto.GetMediaOutput, 0, len(movieIDs))
	for _, movieID := range movieIDs {
		// OPTIMIZATION: Batch fetch movies by IDs to reduce number of queries
		movieDTO, derr := uc.getMediaUseCase.Execute(ctx, dto.GetMediaInput{MediaID: movieID})
		if derr != nil {
			log.Error("Failed to get movie by ID in recommendations", derr)
			return dto.GetMediaRecommendationsOutput{}, derr
		}
		dtoMovies = append(dtoMovies, movieDTO)
	}

	return dto.GetMediaRecommendationsOutput{Movies: dtoMovies}, nil
}
