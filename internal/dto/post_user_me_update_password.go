package dto

type PostUserMeUpdatePasswordInput struct {
	AccessToken     string `json:"access_token" validate:"required"`
	CurrentPassword string `json:"current_password" validate:"required,min=8,max=64"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=64"`
}

type PostUserMeUpdatePasswordOutput struct{}
