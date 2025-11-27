package dto

import "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"

type GetMediaEpisodesInput struct {
	MediaID uint `json:"media_id" validate:"required"`
}

type GetMediaEpisodeOutput struct {
	entity.Episode
	Media  GetMediaOutput `json:"media"`
}

type GetMediaEpisodesOutput struct {
	Episodes []GetMediaEpisodeOutput `json:"episodes"`
}