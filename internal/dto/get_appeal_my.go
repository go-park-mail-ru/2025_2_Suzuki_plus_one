package dto

//go:generate easyjson -all $GOFILE

type GetAppealMyInput struct {
	AccessToken string `json:"access_token"` // Access token from Authorization header
}

type GetAppealMyOutput struct {
	Appeals []GetAppealOutput `json:"appeals"`
}
