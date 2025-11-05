package usecase

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
)

//go:generate mockgen -source=contract.go -destination=./mocks/contract_mock.go -package=mocks
type (
	// Postgres
	Repository interface {
		Connect() error
		Close() error
	}
	MediaRepository interface {
		// Return total count of media items of specific type
		GetMediaCount(ctx context.Context, media_type string) (int, error)

		GetMediaByID(ctx context.Context, media_id uint) (*entity.Media, error)
		// Get genres related to specific media ID
		GetMediaGenres(ctx context.Context, media_id uint) ([]entity.Genre, error)
		// Get posters S3 keys related to specific media ID
		GetMediaPostersKeys(ctx context.Context, media_id uint) ([]entity.S3Key, error)
		// Get media s3 key for watching
		GetMediaWatchKey(ctx context.Context, media_id uint) (*entity.S3Key, error)

		// Get random media IDs for recommendations
		GetMediaRandomIds(ctx context.Context, limit uint, offset uint, media_type string) ([]uint, error)
	}

	ActorRepository interface {
		GetActorByID(ctx context.Context, actorID uint) (*entity.Actor, error)
		// Get actors related to specific media ID
		GetActorsByMediaID(ctx context.Context, media_id uint) ([]entity.Actor, error)
		GetMediasByActorID(ctx context.Context, actorID uint) ([]entity.Media, error)
		GetActorImageS3(ctx context.Context, actorID uint) ([]entity.S3Key, error)
	}

	UserRepository interface {
		GetUserByID(ctx context.Context, userID uint) (*entity.User, error)
		GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
		GetUserAvatarKey(ctx context.Context, userID uint) (*entity.S3Key, error)
		CreateUser(ctx context.Context, user entity.User) (uint, error)
	}

	TokenRepository interface {
		AddNewRefreshToken(ctx context.Context, userID uint, refreshToken string, expiresAt time.Time) error
		GetRefreshTokensForUser(ctx context.Context, userID uint) ([]entity.RefreshToken, error)
		RemoveRefreshToken(ctx context.Context, userID uint, refreshToken string) error
	}

	// Minio
	S3 interface{}

	ObjectRepository interface {
		GetObject(ctx context.Context, bucketName string, key string, expiration time.Duration) (*entity.Object, error)
		GetPublicObject(ctx context.Context, bucketName string, key string) (*entity.Object, error)
	}

	// Redis
	Cache interface{}

	SessionRepository interface {
		AddSession(ctx context.Context, userID uint, accessToken string, expiration time.Duration) error
		DeleteSession(ctx context.Context, userID uint, accessToken string) error
		GetUserIDByToken(ctx context.Context, accessToken string) (uint, error)
	}
)
