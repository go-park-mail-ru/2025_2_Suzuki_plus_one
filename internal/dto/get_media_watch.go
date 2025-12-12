package dto

type GetMediaWatchInput struct {
	AccessToken string `json:"access_token"`
	MediaID     uint   `json:"media_id" validate:"required,min=1"`
}

type GetMediaWatchOutput struct {
	URL string `json:"url"`
}
