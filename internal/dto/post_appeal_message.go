package dto

type PostAppealMessageInput struct {
	AppealId    uint   `json:"appeal_id"`    // ID of the appeal to which the message is being added
	AccessToken string `json:"access_token"` // Access token from Authorization header
	Message     string `json:"message"`      // Content of the message being added to the appeal
}

type PostAppealMessageOutput struct {
	ID uint `json:"id"`
}
