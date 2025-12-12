package aws

import (
	"context"
	"time"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

func (sss *AWSS3) GeneratePresignedURL(ctx context.Context, bucketName string, objectName string, expiration time.Duration) (*entity.URL, error) {
	log := logger.LoggerWithKey(sss.logger, ctx, common.ContextKeyRequestID)

	log.Debug("GeneratePresignedURL called",
		log.ToString("bucketName", bucketName),
		log.ToString("objectName", objectName),
	)

	// Map logical bucket name to actual S3 bucket name
	mappedBucketName, ok := BucketMap[bucketName]
	if !ok || mappedBucketName == "" {
		log.Error("Invalid or missing bucket mapping for: " + bucketName)
		return nil, fmt.Errorf("invalid or missing bucket mapping for: %s", bucketName)
	}

	// Check if the object exists
	_, err := sss.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(mappedBucketName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		log.Error("Failed to stat object: " + err.Error())
		return nil, err
	}

	// Generate presigned URL
	presignClient := s3.NewPresignClient(sss.client)
	presignedReq, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(mappedBucketName),
		Key:    aws.String(objectName),
	}, s3.WithPresignExpires(expiration))
	if err != nil {
		log.Error("Failed to generate presigned URL: " + err.Error())
		return nil, err
	}

	log.Info("Presigned URL generated successfully",
		log.ToString("bucketName", bucketName),
		log.ToString("mappedBucketName", mappedBucketName),
		log.ToString("objectName", objectName),
		log.ToString("presignedURL", presignedReq.URL),
	)

	return &entity.URL{
		URL: presignedReq.URL,
	}, nil
}

func (sss *AWSS3) GeneratePublicURL(ctx context.Context, bucketName string, objectName string) (*entity.URL, error) {
	log := logger.LoggerWithKey(sss.logger, ctx, common.ContextKeyRequestID)

	log.Debug("GeneratePublicURL called",
		log.ToString("bucketName", bucketName),
		log.ToString("objectName", objectName),
	)
	// Map logical bucket name to actual S3 bucket name
	mappedBucketName, ok := BucketMap[bucketName]
	if !ok || mappedBucketName == "" {
		log.Error("Invalid or missing bucket mapping for: " + bucketName)
		return nil, fmt.Errorf("invalid or missing bucket mapping for: %s", bucketName)
	}

	// https://cloud.ru/docs/s3e/ug/topics/guides__public-access-for-bucket?source-platform=Evolution
	// Example public URL format:
	// 		https://global.s3.cloud.ru/bucketName/objectName

	publicURL := sss.publicURL + "/" + mappedBucketName + "/" + objectName

	log.Info("Public URL generated successfully",
		log.ToString("bucketName", mappedBucketName),
		log.ToString("objectName", objectName),
		log.ToString("publicURL", publicURL),
	)

	return &entity.URL{
		URL: publicURL,
	}, nil
}
