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
		// Get trailers S3 keys related to specific media ID
		GetMediaTrailersKeys(ctx context.Context, media_id uint) ([]entity.S3Key, error)
		// Get media s3 key for watching
		GetMediaWatchKey(ctx context.Context, media_id uint) (*entity.S3Key, error)

		// Get random media IDs for recommendations
		GetMediaSortedByName(ctx context.Context, limit uint, offset uint, media_type string) ([]uint, error)
	}

	ActorRepository interface {
		GetActorByID(ctx context.Context, actorID uint) (*entity.Actor, error)
		// Get actors related to specific media ID
		GetActorsByMediaID(ctx context.Context, media_id uint) ([]entity.Actor, error)
		GetMediasByActorID(ctx context.Context, actorID uint) ([]entity.Media, error)
		GetActorImageS3(ctx context.Context, actorID uint) ([]entity.S3Key, error)
	}

	UserRepository interface {
		// Create
		CreateUser(ctx context.Context, email string, username string, passwordHash string) (uint, error)

		// Get
		GetUserByID(ctx context.Context, userID uint) (*entity.User, error)
		GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
		GetUserAvatarKey(ctx context.Context, userID uint) (*entity.S3Key, error)

		// Update
		UpdateUser(ctx context.Context,
			userID uint,
			username string,
			email string,
			dateOfBirth time.Time,
			phoneNumber string,
		) (*entity.User, error)
		UpdateUserAvatarKey(ctx context.Context, userID uint, assetImageID uint) error
		UpdateUserPassword(ctx context.Context, userID uint, newHashedPassword string) error
	}

	AppealRepository interface {
		// Get
		GetAppealIDsByUserID(ctx context.Context, userID uint) ([]uint, error)
		GetAppealIDsAll(ctx context.Context, tag *string, status *string, limit uint, offset uint) ([]uint, error)
		GetAppealByID(ctx context.Context, appealID uint) (*entity.Appeal, error)

		// Get messages for appeal
		GetAppealMessagesByID(ctx context.Context, appealID uint) ([]entity.AppealMessage, error)

		// Create
		CreateAppeal(ctx context.Context, userID uint, tag string, name string) (uint, error)
		CreateAppealMessage(ctx context.Context, appealID uint, isResponse bool, message string) (uint, error)

		// Update
		UpdateAppealStatus(ctx context.Context, appealID uint, status string) error
	}

	AssetRepository interface {
		// Asset
		CreateAsset(ctx context.Context, asset entity.Asset) (uint, error)
		GetAssetByID(ctx context.Context, assetID uint) (*entity.Asset, error)

		// AssetImage
		CreateAssetImage(ctx context.Context, assetImage entity.AssetImage) (uint, error)
		GetAssetImageByID(ctx context.Context, assetImageID uint) (*entity.AssetImage, error)
	}

	TokenRepository interface {
		AddNewRefreshToken(ctx context.Context, userID uint, refreshToken string, expiresAt time.Time) error
		GetRefreshTokensForUser(ctx context.Context, userID uint) ([]entity.RefreshToken, error)
		RemoveRefreshToken(ctx context.Context, userID uint, refreshToken string) error
	}

	// Minio
	S3 interface{}

	ObjectRepository interface {
		GeneratePublicURL(ctx context.Context, bucketName string, objectName string) (*entity.URL, error)
		GeneratePresignedURL(ctx context.Context, bucketName string, objectName string, expiration time.Duration) (*entity.URL, error)
		UploadObject(ctx context.Context, bucketName string, objectName string, mimeType string, data []byte) (*entity.S3Key, error)
		DeleteObject(ctx context.Context, bucketName, objectName string) error
	}

	// Redis
	Cache interface{}

	SessionRepository interface {
		AddSession(ctx context.Context, userID uint, accessToken string, expiration time.Duration) error
		DeleteSession(ctx context.Context, userID uint, accessToken string) error
		GetUserIDByAccessToken(ctx context.Context, accessToken string) (uint, error)
	}
)
