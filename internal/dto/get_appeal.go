package dto

import "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"

type GetAppealInput struct {
	AccessToken string `json:"access_token"` // Access token from Authorization header
	AppealId    uint   `json:"appeal_id"`    // ID of the appeal to retrieve
}

type GetAppealOutput struct {
	entity.Appeal
}
