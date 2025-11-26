package dto

type GetAppealMyInput struct {
	AccessToken string `json:"access_token"` // Access token from Authorization header
}

type GetAppealMyOutput struct {
	Appeals []GetAppealOutput `json:"appeals"`
}
