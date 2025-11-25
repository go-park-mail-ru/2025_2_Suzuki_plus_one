package dto

type GetMediaMyInput struct {
	IsDislike   bool   `json:"is_dislike"`
	AccessToken string `json:"access_token" validate:"required"`
	Limit       uint   `json:"limit" validate:"min=1,max=100"`
	Offset      uint   `json:"offset" validate:"min=0"`
}

type GetMediaMyOutput struct {
	Medias []GetMediaOutput `json:"medias"`
}
