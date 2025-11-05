package minio

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/minio/minio-go/v7"
)

// GetObject retrieves an object from MinIO and generates a presigned URL for access.
func (m *Minio) GetObject(ctx context.Context, bucketName string, objectName string, expiration time.Duration) (*entity.Object, error) {
	requestID, ok := ctx.Value(common.RequestIDContextKey).(string)
	if !ok {
		m.logger.Warn("GetObject: failed to get requestID from context")
		requestID = "unknown"
	}
	m.logger.Info("GetObject called",
		m.logger.ToString("requestID", requestID),
		m.logger.ToString("bucketName", bucketName),
		m.logger.ToString("objectName", objectName),
	)

	// Check if the object exists
	info, err := m.client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		m.logger.Error("Failed to stat object: " + err.Error())
		return nil, err
	}
	if info.Size == 0 {
		m.logger.Error("Object not found or is empty")
		return nil, fmt.Errorf("object not found or is empty")
	}

	// Generate presigned URL
	reqParams := make(url.Values)
	presignedURL, err := m.client.PresignedGetObject(ctx, bucketName, objectName, expiration, reqParams)
	if err != nil {
		m.logger.Error("Failed to generate presigned URL: " + err.Error())
		return nil, err
	}

	// Replace internal URL with public URL
	parsedURL, err := url.Parse(presignedURL.String())
	if err != nil {
		m.logger.Error("Failed to parse presigned URL: " + err.Error())
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
	requestID, ok := ctx.Value(common.RequestIDContextKey).(string)
	if !ok {
		m.logger.Warn("GetPublicObject: failed to get requestID from context")
		requestID = "unknown"
	}
	m.logger.Info("GetPublicObject called",
		m.logger.ToString("requestID", requestID),
		m.logger.ToString("bucketName", bucketName),
		m.logger.ToString("objectName", objectName),
	)

	// Check if the object exists
	info, err := m.client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		m.logger.Error("Failed to stat object: " + err.Error())
		return nil, err
	}
	if info.Size == 0 {
		m.logger.Error("Object not found or is empty")
		return nil, fmt.Errorf("object not found or is empty")
	}

	// Construct public URL
	objectURL := fmt.Sprintf("%s/%s/%s", m.externalHost, bucketName, objectName)
	return &entity.Object{
		URL: objectURL,
	}, nil
}
