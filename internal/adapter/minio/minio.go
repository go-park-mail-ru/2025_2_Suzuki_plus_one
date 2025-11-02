package minio

import (
	"context"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	logger  logger.Logger
	client  *minio.Client
	context context.Context
}

func NewMinio(logger logger.Logger, endpoint string, login string, password string, useSSL bool) (*Minio, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(login, password, ""),
		Secure: useSSL,
	})
	if err != nil {
		logger.Error("Failed to create Minio client: " + err.Error())
		return nil, err
	}
	return &Minio{
		logger:  logger,
		client:  minioClient,
		context: context.Background(),
	}, nil
}
