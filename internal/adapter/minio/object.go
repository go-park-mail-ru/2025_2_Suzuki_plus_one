package minio

import (
	"bytes"
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/minio/minio-go/v7"
)

func (m *Minio) UploadObject(ctx context.Context, bucketName string, objectName string, mimeType string, data []byte) (*entity.S3Key, error) {
	log := logger.LoggerWithKey(m.logger, ctx, common.ContextKeyRequestID)
	log.Debug("UploadObject called",
		log.ToString("bucketName", bucketName),
		log.ToString("objectName", objectName),
	)

	// Upload the object
	_, err := m.client.PutObject(
		ctx,
		bucketName,
		objectName,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{ContentType: mimeType},
	)
	if err != nil {
		log.Error("Failed to upload object: " + err.Error())
		return &entity.S3Key{}, err
	}

	log.Info("File uploaded to Minio successfully",
		log.ToString("bucketName", bucketName),
		log.ToString("objectName", objectName),
	)
	return &entity.S3Key{
		Key:        objectName,
		BucketName: bucketName,
	}, nil
}

func (m *Minio) DeleteObject(ctx context.Context, bucketName string, objectName string) error {
	log := logger.LoggerWithKey(m.logger, ctx, common.ContextKeyRequestID)
	log.Debug("DeleteObject called",
		log.ToString("bucketName", bucketName),
		log.ToString("objectName", objectName),
	)

	err := m.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		log.Error("Failed to delete object: " + err.Error())
		return err
	}

	log.Info("File deleted from Minio successfully",
		log.ToString("bucketName", bucketName),
		log.ToString("objectName", objectName),
	)
	return nil
}
