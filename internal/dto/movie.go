package dto

import "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"

type Movie struct {
	entity.Media
	Genres  []entity.Genre `json:"genres"`
	Posters []string       `json:"posters"`
}

func NewMovieFromEntity(e *entity.Media, genres []entity.Genre, posters []string) Movie {
	return Movie{
		Media:   *e,
		Genres:  genres,
		Posters: posters,
	}
}
