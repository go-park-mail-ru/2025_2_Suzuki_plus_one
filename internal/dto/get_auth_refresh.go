package dto

//go:generate easyjson -all $GOFILE

type GetAuthRefreshInput struct {
	RefreshToken string `json:"refresh_token" required:"true"`
}

type GetAuthRefreshOutput struct {
	AccessToken string `json:"access_token"`
}
