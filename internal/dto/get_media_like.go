package dto

type GetMediaLikeInput struct {
	AccessToken string `json:"access_token" validate:"required"`
	MediaID     uint   `json:"media_id" validate:"required"`
}

type GetMediaLikeOutput struct {
	Liked     bool `json:"liked"`
	IsDislike bool `json:"is_dislike"`
}
