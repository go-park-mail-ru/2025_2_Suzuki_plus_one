package dto

import "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"

type GetUserMeInput struct {
	AccessToken string `json:"access_token"` // Access token from Authorization header
}

type GetUserMeOutput struct {
	entity.User
}
