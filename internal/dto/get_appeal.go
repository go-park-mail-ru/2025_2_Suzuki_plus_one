package dto

type GetAppealInput struct {
	AppealId    uint   `json:"appeal_id"`    // ID of the appeal to retrieve
	AccessToken string `json:"access_token"` // Access token from Authorization header
}

type GetAppealOutput struct {
	ID        uint   `json:"id"`
	Tag       string `json:"tag"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type AppealMessage struct {
	ID        uint   `json:"id"`
	Sender    string `json:"sender"` // "user" or "support"
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}
