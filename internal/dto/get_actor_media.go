package dto

type GetActorMediaInput struct {
	ActorID uint `json:"actor_id"`
}

type GetActorMediaOutput struct {
	Medias []GetMediaOutput `json:"medias"`
}
