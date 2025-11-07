package entity

type Asset struct {
	ID         uint    `json:"id"`
	S3Key      string  `json:"s3_key"`
	MimeType   string  `json:"mime_type"`
	FileSizeMB float32 `json:"file_size_mb"`
}

type AssetImage struct {
	ID               uint `json:"id"`
	AssetID          uint `json:"asset_id"`
	ResolutionWidth  uint `json:"resolution_width"`
	ResolutionHeight uint `json:"resolution_height"`
}
