package dto

//go:generate easyjson -all $GOFILE
type GetGenreAllInput struct {
}

//go:generate easyjson -all $GOFILE
type GetGenreAllOutput struct {
	Genres []GetGenreOutput `json:"genres"`
}
