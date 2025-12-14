package dto

import "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"

//go:generate easyjson -all $GOFILE

type GetGenreInput struct {
	GenreID uint `json:"genre_id" validate:"required"`
}

type GetGenreOutput struct {
	entity.Genre
}
