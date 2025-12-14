package dto

//go:generate easyjson -all $GOFILE
type PostAuthSignInInput struct {
	Email    string `json:"email" required:"true"`
	Password string `json:"password" required:"true"`
}


//go:generate easyjson -all $GOFILE
type PostAuthSignInOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
