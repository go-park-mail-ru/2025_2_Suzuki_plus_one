package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
)

type GetMoviesUsecase struct {
	movieRepo MovieRepository
}

func NewGetMoviesUsecase(movieRepo MovieRepository) *GetMoviesUsecase {
	return &GetMoviesUsecase{
		movieRepo: movieRepo,
	}
}

func (uc *GetMoviesUsecase) Execute(ctx context.Context, input dto.GetMoviesInput) (dto.GetMoviesOutput, *dto.Error) {
	// Validate input.
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError("usecase/get_movies", entity.ErrGetMoviesParamsInvalid, err.Error())
		return dto.GetMoviesOutput{}, &derr
	}

	// Get movies from repository
	movies, err := uc.movieRepo.GetMovies(ctx, input.Limit, input.Offset)
	if err != nil {
		derr := dto.NewError("adapter/get_movies", err, "")
		return dto.GetMoviesOutput{}, &derr
	}

	// Convert []entity.Movie to []dto.Movie
	dtoMovies := make([]dto.Movie, 0, len(movies))
	for _, m := range movies {
		dtoMovies = append(dtoMovies, dto.NewMovieFromEntity(m))
	}

	return dto.GetMoviesOutput{Movies: dtoMovies}, nil
}
