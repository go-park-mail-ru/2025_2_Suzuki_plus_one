package dto

type GetAppealAllInput struct {
	AccessToken string `json:"access_token"` // Access token from Authorization header
}

type GetAppealAllOutput struct {
	Appeals GetAppealOutput `json:"appeals"`
}
