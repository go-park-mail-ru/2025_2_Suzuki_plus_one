package dto

type PutAppealResolveInput struct {
	AppealId    uint   `json:"appeal_id"`    // ID of the appeal to resolve
	AccessToken string `json:"access_token"` // Access token from Authorization header
}

type PutAppealResolveOutput struct {
	ID     uint   `json:"id"`
	Status string `json:"status"`
}
