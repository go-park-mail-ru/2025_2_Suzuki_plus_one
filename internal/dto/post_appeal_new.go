package dto

//go:generate easyjson -all $GOFILE

type PostAppealNewInput struct {
	AccessToken string `json:"access_token"` // Access token from Authorization header
	Tag         string `json:"tag"`          // Tag or category of the appeal
	Message     string `json:"message"`      // Initial message of the appeal
	Name        string `json:"name"`         // Name of the appeal
}

type PostAppealNewOutput struct {
	ID uint `json:"id"`
}
