package app

// Postgres
type DatabaseRepository interface {
	Connect() error
	Close() error
}

// Minio
type S3 interface{}

// Redis
type Cache interface {
	Close() error
}

// gRPC Auth service
type Service interface {
	Connect() error
	Close()
}
