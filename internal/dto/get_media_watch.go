package dto

//go:generate easyjson -all $GOFILE

type GetMediaWatchInput struct {
	AccessToken string `json:"access_token" validate:"required"`
	MediaID     uint   `json:"media_id" validate:"required,min=1"`
}

type GetMediaWatchOutput struct {
	URL string `json:"url"`
}
