package dto

type GetMovieRecommendationsInput struct {
	Limit  uint `json:"limit" validate:"gte=0" required:"true"`
	Offset uint `json:"offset" validate:"gte=0"`
}

type GetMovieRecommendationsOutput struct {
	Movies []GetMediaOutput `json:"movies"`
}
