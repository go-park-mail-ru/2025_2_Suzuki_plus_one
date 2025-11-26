package minio

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	logger       logger.Logger
	client       *minio.Client
	context      context.Context
	secure       bool
	internalHost string
	externalHost string
}

func NewMinio(logger logger.Logger, internalHost string, externalHost string, login string, password string, useSSL bool) (*Minio, error) {
	minioClient, err := minio.New(internalHost, &minio.Options{
		Creds:  credentials.NewStaticV4(login, password, ""),
		Secure: useSSL,
	})
	if err != nil {
		logger.Error("Failed to create Minio client: " + err.Error())
		return nil, err
	}

	buckets, err := minioClient.ListBuckets(context.Background())
	if err != nil {
		logger.Fatal("Failed to connect to MinIO: ", err)
	}
	logger.Info("Connected to MinIO successfully. Available buckets:")
	for _, bucket := range buckets {
		logger.Info(" - " + bucket.Name)
	}

	return &Minio{
		logger:       logger,
		client:       minioClient,
		context:      context.Background(),
		internalHost: internalHost,
		externalHost: externalHost,
		secure:       useSSL,
	}, nil
}
