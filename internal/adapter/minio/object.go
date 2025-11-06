package minio

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/minio/minio-go/v7"
)

// GetObject retrieves an object from MinIO and generates a presigned URL for access.
func (m *Minio) GetObject(ctx context.Context, bucketName string, objectName string, expiration time.Duration) (*entity.Object, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(m.logger, ctx, common.ContexKeyRequestID)

	log.Info("GetObject called",
		log.ToString("bucketName", bucketName),
		log.ToString("objectName", objectName),
	)

	// Check if the object exists
	info, err := m.client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		log.Error("Failed to stat object: " + err.Error())
		return nil, err
	}
	if info.Size == 0 {
		log.Error("Object not found or is empty")
		return nil, fmt.Errorf("object not found or is empty")
	}

	// Generate presigned URL
	reqParams := make(url.Values)
	presignedURL, err := m.client.PresignedGetObject(ctx, bucketName, objectName, expiration, reqParams)
	if err != nil {
		log.Error("Failed to generate presigned URL: " + err.Error())
		return nil, err
	}

	// Replace internal URL with public URL
	parsedURL, err := url.Parse(presignedURL.String())
	if err != nil {
		log.Error("Failed to parse presigned URL: " + err.Error())
		return nil, err
	}

	// Use the host from public URL
	publicParsedURL, _ := url.Parse(m.externalHost)
	parsedURL.Host = publicParsedURL.Host
	parsedURL.Scheme = publicParsedURL.Scheme

	return &entity.Object{
		URL: parsedURL.String(),
	}, nil
}

func (m *Minio) GetPublicObject(ctx context.Context, bucketName string, objectName string) (*entity.Object, error) {
	log := logger.LoggerWithKey(m.logger, ctx, common.ContexKeyRequestID)
	log.Debug("GetPublicObject called",
		log.ToString("bucketName", bucketName),
		log.ToString("objectName", objectName),
	)

	// Check if the object exists
	info, err := m.client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		log.Error("Failed to stat object: " + err.Error())
		return nil, err
	}
	if info.Size == 0 {
		log.Error("Object not found or is empty")
		return nil, fmt.Errorf("object not found or is empty")
	}

	// Construct public URL
	objectURL := fmt.Sprintf("%s/%s/%s", m.externalHost, bucketName, objectName)
	return &entity.Object{
		URL: objectURL,
	}, nil
}
