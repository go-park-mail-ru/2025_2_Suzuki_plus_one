package dto

type GetObjectInput struct {
	Key        string `json:"key" validate:"required"`
	BucketName string `json:"bucket_name" validate:"required,oneof=avatars posters trailers media"`
}

type GetObjectOutput struct {
	URL string `json:"url"`
}
