package usecase

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type GetMovieRecommendationsUsecase struct {
	logger    logger.Logger
	movieRepo MovieRepository
}

func NewGetMovieRecommendationsUsecase(
	logger logger.Logger,
	movieRepo MovieRepository,
) *GetMovieRecommendationsUsecase {
	return &GetMovieRecommendationsUsecase{
		logger:    logger,
		movieRepo: movieRepo,
	}
}

func (uc *GetMovieRecommendationsUsecase) Execute(
	ctx context.Context,
	input dto.GetMovieRecommendationsInput,
) (
	dto.GetMovieRecommendationsOutput,
	*dto.Error,
) {
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
		return dto.GetMovieRecommendationsOutput{}, &derr
	}

	// Calculate recommendation movie IDs (dummy logic for now)
	movieIDs := getRecommendedMovieIDs(input.Limit, input.Offset, uint(movie_number))

	// Convert []entity.Movie to []dto.Movie
	dtoMovies := make([]dto.Movie, 0, len(movieIDs))
	for _, movieID := range movieIDs {
		// OPTIMIZATION: Batch fetch movies by IDs to reduce number of queries
		movie, err := uc.movieRepo.GetMedia(ctx, movieID)
		if err != nil {
			derr := dto.NewError("adapter/get_movie_recommendations", err, "Can't get movie by ID")
			return dto.GetMovieRecommendationsOutput{}, &derr
		}

		// Fetch genres and posters for current movie
		genres, err := uc.movieRepo.GetMediaGenres(ctx, movieID)
		if err != nil {
			derr := dto.NewError("adapter/get_movie_recommendations", err, "Can't get movie genres")
			return dto.GetMovieRecommendationsOutput{}, &derr
		}
		posters, err := uc.movieRepo.GetMediaPostersLinks(ctx, movieID)
		if err != nil {
			derr := dto.NewError("adapter/get_movie_recommendations", err, "Can't get movie posters")
			return dto.GetMovieRecommendationsOutput{}, &derr
		}

		dtoMovies = append(dtoMovies, dto.NewMovieFromEntity(movie, genres, posters))
	}

	return dto.GetMovieRecommendationsOutput{Movies: dtoMovies}, nil
}

// Dummy function to generate random recommended movie IDs
func getRecommendedMovieIDs(limit, offset, totalMovies uint) []uint {
	movieIDs := make([]uint, 0, limit)

	// Simple random selection logic
	// OPTIMIZATION: Use better randomization algorithm to avoid duplicates
	// Or just switch to proper recommendation algorithm
	random_source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(random_source)
	used := make(map[uint]struct{})
	for uint(len(movieIDs)) < limit && totalMovies > 0 {
		id := rng.Uint32()%uint32(totalMovies) + 1
		if _, exists := used[uint(id)]; !exists {
			movieIDs = append(movieIDs, uint(id))
			used[uint(id)] = struct{}{}
		}
	}
	return movieIDs
}
