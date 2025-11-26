package dto

import "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"

type GetAppealMessageInput struct {
	AppealID    uint   `json:"appeal_id"`    // ID of the appeal to which the message is being returned
	AccessToken string `json:"access_token"` // Access token from Authorization header
}

type GetAppealMessageOutput struct {
	Messages []AppealMessage `json:"messages"`
}

// We overload AppealMessage to change CreatedAt type to JSONDateTime format
type AppealMessage struct {
	entity.AppealMessage
	CreatedAt JSONDateTime `json:"timestamp"`
}
