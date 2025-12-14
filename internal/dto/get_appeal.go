package dto

import "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"

//go:generate easyjson -all $GOFILE

type GetAppealInput struct {
	AccessToken string `json:"access_token"` // Access token from Authorization header
	AppealId    uint   `json:"appeal_id"`    // ID of the appeal to retrieve
}

type GetAppealOutput struct {
	Appeal
}

type Appeal struct {
	entity.Appeal
	CreatedAt JSONDateTime `json:"created_at"`
	UpdatedAt JSONDateTime `json:"updated_at"`
}
