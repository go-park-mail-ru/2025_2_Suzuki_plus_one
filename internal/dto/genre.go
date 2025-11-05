package dto

import "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"

type GenreOutput struct {
	entity.Genre
}

func NewGenreOutputFromEntity(e entity.Genre) GenreOutput {
	return GenreOutput{
		Genre: e,
	}
}
