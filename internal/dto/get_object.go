package dto

//go:generate easyjson -all $GOFILE
type GetObjectInput struct {
	Key        string `json:"key" validate:"required"`
	BucketName string `json:"bucket_name" validate:"required,oneof=actors avatars posters trailers medias"`
}

//go:generate easyjson -all $GOFILE
type GetObjectOutput struct {
	URL string `json:"url"`
}
