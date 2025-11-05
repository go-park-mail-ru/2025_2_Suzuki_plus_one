package entity

type Object struct {
	URL string `json:"url"`
}

type S3Key struct {
	Key        string
	BucketName string
}

func (k S3Key) GetPath() string {
	return k.BucketName + "/" + k.Key
}
