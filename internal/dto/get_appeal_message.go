package dto

type GetAppealMessageInput struct {
	AppealId    uint   `json:"appeal_id"`    // ID of the appeal to which the message is being returned
	AccessToken string `json:"access_token"` // Access token from Authorization header
}

type GetAppealMessageOutput struct {
	Messages []AppealMessage `json:"messages"`
}
