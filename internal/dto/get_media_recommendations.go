package dto

//go:generate easyjson -all $GOFILE

type GetMediaRecommendationsInput struct {
	Limit    uint   `json:"limit" validate:"gte=0" required:"true"`
	Offset   uint   `json:"offset" validate:"gte=0"`
	Type     string `json:"type" validate:"oneof=movie series" required:"true"`
	GenreIDs []uint `json:"genre_ids,omitempty" validate:"dive,gte=0"`
}

type GetMediaRecommendationsOutput struct {
	Movies []GetMediaOutput `json:"movies"`
}
