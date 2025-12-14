package dto

//go:generate easyjson -all $GOFILE

type GetActorMediaInput struct {
	ActorID uint `json:"actor_id"`
}

type GetActorMediaOutput struct {
	Medias []GetMediaOutput `json:"medias"`
}
