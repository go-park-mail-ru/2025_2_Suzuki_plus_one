package dto

import "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"

// Embed entity.Movie to extend it if needed
type Movie struct {
	entity.Movie
}

func NewTestMovie(id string) Movie {
	return Movie{
		Movie: entity.NewTestMovie(id),
	}
}

func NewMovieFromEntity(e entity.Movie) Movie {
	return Movie{
		Movie: e,
	}
}
