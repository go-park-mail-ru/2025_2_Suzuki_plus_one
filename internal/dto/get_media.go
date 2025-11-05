package dto

import "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"

type GetMediaInput struct {
	MediaID uint `json:"media_id"`
}

type GetMediaOutput struct {
	entity.Media
	Genres  []GenreOutput    `json:"genres"`
	Posters []string         `json:"posters"`
	Actors  []GetActorOutput `json:"actors"`
}
