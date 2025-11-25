package dto

type DeleteMediaLikeInput struct {
	AccessToken string `json:"access_token" validate:"required"`
	MediaID     int    `json:"media_id" validate:"required"`
}

type DeleteMediaLikeOutput struct {
}
