package dto

type GetSearchInput struct {
	Query  string `json:"query" validate:"required,min=1,max=100"`
	Limit  uint   `json:"limit" validate:"omitempty,min=1,max=100"`
	Offset uint   `json:"offset" validate:"omitempty,min=0"`
	Type   string `json:"type" validate:"omitempty,oneof=media actor any"`
}

type GetSearchOutput struct {
	Medias []GetMediaOutput `json:"medias"`
	Actors []GetActorOutput `json:"actors"`
}