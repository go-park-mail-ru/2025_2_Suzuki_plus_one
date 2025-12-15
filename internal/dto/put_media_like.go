package dto

//go:generate easyjson -all $GOFILE
type PutMediaLikeInput struct {
	AccessToken string `json:"access_token" validate:"required"`
	MediaID     uint   `json:"media_id" validate:"required"`
	IsDislike   bool   `json:"is_dislike"`
}

type PutMediaLikeOutput struct {
	Liked     bool `json:"liked"`
	IsDislike bool `json:"is_dislike"`
}
