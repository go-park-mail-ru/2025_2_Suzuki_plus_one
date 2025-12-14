package dto

import "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"

//go:generate easyjson -all $GOFILE
type GetActorInput struct {
	ActorID uint `json:"actor_id"`
}

//go:generate easyjson -all $GOFILE
type GetActorOutput struct {
	entity.Actor
	ImageURLs []string `json:"image_urls"`
}
