package aws

import (
	"bytes"
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

func (sss *AWSS3) UploadObject(ctx context.Context, bucketName string, objectName string, mimeType string, data []byte) (*entity.S3Key, error) {
	log := logger.LoggerWithKey(sss.logger, ctx, common.ContextKeyRequestID)
	log.Debug("UploadObject called",
		log.ToString("bucketName", bucketName),
		log.ToString("objectName", objectName),
	)
	// Map logical bucket name to actual S3 bucket name
	bucketName = BucketMap[bucketName]

	// Upload the object
	_, err := sss.client.PutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket:      aws.String(bucketName),
			Key:         aws.String(objectName),
			Body:        bytes.NewReader(data),
			ContentType: aws.String(mimeType),
		},
	)
	if err != nil {
		log.Error("Failed to upload object: " + err.Error())
		return &entity.S3Key{}, err
	}

	log.Info("File uploaded to AWS S3 successfully",
		log.ToString("bucketName", bucketName),
		log.ToString("objectName", objectName),
	)

	// Wait until the object is available
	err = s3.NewObjectExistsWaiter(sss.client).Wait(
		ctx, &s3.HeadObjectInput{Bucket: aws.String(bucketName), Key: aws.String(objectName)}, time.Minute)
	if err != nil {
		log.Error("Failed attempt to wait for object %s to exist.\n", objectName)
	}

	return &entity.S3Key{
		Key:        objectName,
		BucketName: bucketName,
	}, nil
}

func (sss *AWSS3) DeleteObject(ctx context.Context, bucketName string, objectName string) error {
	log := logger.LoggerWithKey(sss.logger, ctx, common.ContextKeyRequestID)
	log.Debug("DeleteObject called",
		log.ToString("bucketName", bucketName),
		log.ToString("objectName", objectName),
	)
	// Map logical bucket name to actual S3 bucket name
	bucketName = BucketMap[bucketName]

	_, err := sss.client.DeleteObject(
		ctx,
		&s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectName),
		},
	)
	if err != nil {
		log.Error("Failed to delete object: " + err.Error())
		return err
	}

	log.Info("File deleted from AWS S3 successfully",
		log.ToString("bucketName", bucketName),
		log.ToString("objectName", objectName),
	)

	return nil
}
