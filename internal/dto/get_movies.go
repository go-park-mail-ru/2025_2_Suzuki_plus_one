package dto

type GetMoviesInput struct {
	Limit  uint `json:"limit" validate:"gte=0"`
	Offset uint `json:"offset" validate:"gte=0"`
}

type GetMoviesOutput struct {
	Movies []Movie `json:"movies"`
}
