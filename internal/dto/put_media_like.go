package dto

type PutMediaLikeInput struct {
	AccessToken string `json:"access_token" validate:"required"`
	MediaID     uint    `json:"media_id" validate:"required"`
}

type PutMediaLikeOutput struct {
	Liked bool `json:"liked"`
	IsDislike bool `json:"is_dislike"`
}
