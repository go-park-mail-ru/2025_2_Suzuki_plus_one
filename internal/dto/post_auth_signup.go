package dto

//go:generate easyjson -all $GOFILE

type PostAuthSignUpInput struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=8"`
}

type PostAuthSignUpOutput struct {
	UserID       int64  `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
