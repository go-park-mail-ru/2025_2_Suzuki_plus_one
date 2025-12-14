package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

var BucketMap = map[string]string{
	"avatars":  "ava",
	"actors":   "actors",
	"posters":  "posters",
	"trailers": "trailers",
	"medias":   "medias",
}

type AWSS3 struct {
	logger    logger.Logger
	client    *s3.Client
	publicURL string
}

// See AWS SDK for Go V2 documentation:
// https://docs.aws.amazon.com/code-library/latest/ug/go_2_s3_code_examples.html#basics

func NewAWSS3(logger logger.Logger, publicURL string) (*AWSS3, error) {
	if logger == nil {
		panic("NewAWSS3: logger cannot be nil")
	}

	if publicURL == "" {
		panic("NewAWSS3: public url cannot be empty")
	}
	// Remove trailing slash if exists
	if publicURL[len(publicURL)-1] == '/' {
		publicURL = publicURL[:len(publicURL)-1]
	}

	// Require aws env variables to be set
	// - AWS_ACCESS_KEY_ID
	// - AWS_SECRET_ACCESS_KEY
	// - AWS_DEFAULT_REGION
	// - AWS_S3_ENDPOINT

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		logger.Error("Failed to load AWS config: " + err.Error())
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)

	// Try to list buckets to ensure S3 client is accessible
	buckets, err := s3Client.ListBuckets(context.Background(), &s3.ListBucketsInput{})
	if err != nil {
		logger.Error("Failed to access S3: " + err.Error())
		return nil, err
	}
	logger.Info("Connected to AWS S3 successfully", logger.ToAny("buckets", buckets.Buckets))

	return &AWSS3{
		logger:    logger,
		client:    s3Client,
		publicURL: publicURL,
	}, nil
}
