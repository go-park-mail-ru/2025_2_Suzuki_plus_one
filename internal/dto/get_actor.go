package dto

import "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"

type GetActorInput struct {
	ActorID uint `json:"actor_id"`
}

type GetActorOutput struct {
	entity.Actor
	Medias    []GetMediaOutput `json:"medias"`
	ImageURLs []string         `json:"image_urls"`
}
