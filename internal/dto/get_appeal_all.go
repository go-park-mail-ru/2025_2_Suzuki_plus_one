package dto

type GetAppealAllInput struct {
	Tag    string `query:"tag" validate:"omitempty,oneof=bug feature other"`
	Status string `query:"status" validate:"omitempty,oneof=open closed in_progress"`
	Limit  uint   `query:"limit" validate:"omitempty,min=1,max=100"`
	Offset uint   `query:"offset" validate:"omitempty,min=0"`
}

type GetAppealAllOutput struct {
	Appeals []GetAppealOutput `json:"appeals"`
}
