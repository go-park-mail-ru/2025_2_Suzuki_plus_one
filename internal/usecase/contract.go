package usecase

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
)

//go:generate mockgen -source=contract.go -destination=./contract_mock.go -package=usecase
type (
	// Postgres
	Repository interface {
		Connect() error
		Close() error
	}
	MovieRepository interface {
		GetMediaCount(ctx context.Context, media_type string) (int, error)
		GetMedia(ctx context.Context, media_id uint) (*entity.Media, error)
		GetMediaGenres(ctx context.Context, media_id uint) ([]entity.Genre, error)
		GetMediaPostersLinks(ctx context.Context, media_id uint) ([]string, error)
	}

	UserRepository interface {
		GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
		GetUserAvatarURL(ctx context.Context, userID uint) (string, error)
	}

	TokenRepository interface {
		AddNewRefreshToken(ctx context.Context, userID uint, refreshToken string, expiresAt time.Time) error
		GetRefreshTokensForUser(ctx context.Context, userID uint) ([]entity.RefreshToken, error)
	}

	// Minio
	S3 interface{}

	ObjectRepository interface {
		GetObject(ctx context.Context, key string, bucketName string, expiration time.Duration) (*entity.Object, error)
	}

	// Redis
	Cache interface{}

	SessionRepository interface {
		AddSession(ctx context.Context, userID uint, accessToken string, expiration time.Duration) error
	}
)
