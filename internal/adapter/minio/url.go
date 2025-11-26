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

// Constructs a public URL for the given object stored in MinIO.
func (m *Minio) GeneratePublicURL(ctx context.Context, bucketName string, objectName string) (*entity.URL, error) {
	log := logger.LoggerWithKey(m.logger, ctx, common.ContextKeyRequestID)
	log.Debug("GeneratePublicURL called",
		log.ToString("bucketName", bucketName),
		log.ToString("objectName", objectName),
	)

	// Check if the object exists
	info, err := m.client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		log.Error("Failed to stat object: "+err.Error(),
			log.ToString("bucketName", bucketName),
			log.ToString("objectName", objectName),
		)
		return nil, err
	}
	if info.Size == 0 {
		log.Error("Object not found or is empty")
		return nil, fmt.Errorf("object not found or is empty")
	}

	// Construct public URL
	objectURL := fmt.Sprintf("%s/%s/%s", m.externalHost, bucketName, objectName)
	return &entity.URL{
		URL: objectURL,
	}, nil
}

// Generates a temporary URL for accessing the given object stored in MinIO.
// TODO: add access rules if needed
func (m *Minio) GeneratePresignedURL(ctx context.Context, bucketName string, objectName string, expiration time.Duration) (*entity.URL, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(m.logger, ctx, common.ContextKeyRequestID)

	log.Info("GeneratePresignedURL called",
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

	return &entity.URL{
		URL: parsedURL.String(),
	}, nil
}
