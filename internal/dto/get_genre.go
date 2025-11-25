package dto

import "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"

type GetGenreInput struct {
	GenreID     uint `json:"genre_id" validate:"required"`
	MediaLimit  uint `json:"media_limit" validate:"required" default:"10"`
	MediaOffset uint `json:"media_offset" default:"0"`
}

type GetGenreOutput struct {
	entity.Genre
	Medias []GetMediaOutput `json:"medias"`
}
