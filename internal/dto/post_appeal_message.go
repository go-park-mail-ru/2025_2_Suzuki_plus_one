package dto

//go:generate easyjson -all $GOFILE
type PostAppealMessageInput struct {
	AccessToken string `json:"access_token"` // Access token from Authorization header
	AppealID    uint   `json:"appeal_id"`    // ID of the appeal to which the message is being added
	Message     string `json:"message"`      // Content of the message being added to the appeal
}


//go:generate easyjson -all $GOFILE
type PostAppealMessageOutput struct {
	ID uint `json:"id"`
}
