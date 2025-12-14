package dto

//go:generate easyjson -all $GOFILE

type PutAppealResolveInput struct {
	AccessToken string `json:"access_token"` // Access token from Authorization header
	AppealId    uint   `json:"appeal_id"`    // ID of the appeal to resolve
}

type PutAppealResolveOutput struct {
	ID      uint   `json:"id"`
	Message string `json:"message"`
}
