package dto

import "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"


type GetMoviesInput struct {
	Limit  uint `json:"limit" validate:"gte=0"`
	Offset uint `json:"offset" validate:"gte=0"`
}

type GetMoviesOutput struct {
	Movies []entity.Movie `json:"movies"`
}