package entity

type URL struct {
	URL string `json:"url"`
}

type S3Key struct {
	Key        string
	BucketName string
}

// Returns the full path of the S3 object in the format "bucketName/key"
func (k S3Key) GetPath() string {
	return k.BucketName + "/" + k.Key
}
