package dto

//go:generate easyjson -all $GOFILE

type GetGenreMediaInput struct {
	GenreID uint `json:"genre_id" validate:"required"`
	Limit   uint `json:"limit" validate:"gte=0,max=100"`
	Offset  uint `json:"offset" validate:"gte=0"`
}

type GetGenreMediaOutput struct {
	Medias []GetMediaOutput `json:"medias"`
}
