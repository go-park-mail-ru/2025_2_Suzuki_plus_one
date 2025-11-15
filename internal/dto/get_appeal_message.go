package dto

import "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"

type GetAppealMessageInput struct {
	AppealId    uint   `json:"appeal_id"`    // ID of the appeal to which the message is being returned
	AccessToken string `json:"access_token"` // Access token from Authorization header
}

type GetAppealMessageOutput struct {
	Messages []entity.AppealMessage `json:"messages"`
}
