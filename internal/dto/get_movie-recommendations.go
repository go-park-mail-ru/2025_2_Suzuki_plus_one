package dto

type GetMovieRecommendationsInput struct {
	Limit  uint `json:"limit" validate:"gte=0"`
	Offset uint `json:"offset" validate:"gte=0"`
}

type GetMovieRecommendationsOutput struct {
	Movies []Movie `json:"movies"`
}
