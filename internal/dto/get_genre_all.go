package dto

type GetGenreAllInput struct {
}

type GetGenreAllOutput struct {
	Genres []GetGenreOutput `json:"genres"`
}
