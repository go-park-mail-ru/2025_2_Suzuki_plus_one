package dto

//go:generate easyjson -all $GOFILE

type GetMediaActorInput struct {
	MediaID uint `json:"media_id"`
}

type GetMediaActorOutput struct {
	Actors []GetActorOutput `json:"actors"`
}
