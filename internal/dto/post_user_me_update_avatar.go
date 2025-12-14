package dto

//go:generate easyjson -all $GOFILE

type PostUserMeUpdateAvatarInput struct {
	AccessToken string  `json:"access_token" required:"true"`
	Bytes       []byte  `json:"avatar" required:"true"`
	MimeFormat  string  `json:"mime_format" required:"true" validate:"oneof=image/jpeg image/png image/jpg"`
	FileSizeMB  float32 `json:"file_size_mb" required:"true" validate:"max=10"`
}

type PostUserMeUpdateAvatarOutput struct {
	URL string `json:"url"`
}
