package dto

type GetAuthSignOutInput struct {
	AccessToken  string `json:"access_token" required:"true"`
	RefreshToken string `json:"refresh_token" required:"true"`
}

type GetAuthSignOutOutput struct {
	// + example:"refresh_token=; HttpOnly; Strict; Max-Age=0"
}
