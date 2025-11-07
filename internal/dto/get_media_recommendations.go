package dto

type GetMediaRecommendationsInput struct {
	Limit  uint   `json:"limit" validate:"gte=0" required:"true"`
	Offset uint   `json:"offset" validate:"gte=0"`
	Type   string `json:"type" validate:"oneof=movie series" required:"true"`
}

type GetMediaRecommendationsOutput struct {
	Movies []GetMediaOutput `json:"movies"`
}
