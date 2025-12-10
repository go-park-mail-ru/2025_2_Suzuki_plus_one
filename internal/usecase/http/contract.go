package http

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	yoopayment "github.com/rvinnie/yookassa-sdk-go/yookassa/payment"
)

//go:generate mockgen -source=contract.go -destination=./mocks/contract_mock.go -package=mocks
type (
	// Postgres
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
		GetMediaSortedByName(ctx context.Context, limit uint, offset uint, media_type string, media_prefered_genres []uint) ([]uint, error)

		// Get media IDs by like status
		GetMediaIDsByLikeStatus(ctx context.Context, userID uint, isDislike bool, limit uint, offset uint) ([]uint, error)

		// Get media IDs related to specific genre ID
		GetMediasByGenreID(ctx context.Context, limit uint, offset uint, genreID uint) ([]uint, error)

		// Get episodes (medias with type episode) related to specific media ID (with type series)
		GetEpisodesByMediaID(ctx context.Context, media_id uint) ([]entity.Episode, error)
	}

	GenreRepository interface {
		// Get genre by ID
		GetGenreByID(ctx context.Context, genreID uint) (*entity.Genre, error)

		// Get all genres Ids
		GetAllGenreIDs(ctx context.Context) ([]uint, error)
	}

	LikeRepository interface {
		// Get
		GetLike(ctx context.Context, userID uint, mediaID uint) (exists bool, ent entity.Like, err error)
		// Create
		ToggleLike(ctx context.Context, userID uint, mediaID uint) (isDislike bool, err error)
		// Delete
		DeleteLike(ctx context.Context, userID uint, mediaID uint) error
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
		GetUserSubscriptionStatus(ctx context.Context, userID uint) (string, error)

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
		UpdateUserSubscriptionStatus(ctx context.Context, userID uint, status string) error
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

	// S3
	ObjectRepository interface {
		GeneratePublicURL(ctx context.Context, bucketName string, objectName string) (*entity.URL, error)
		GeneratePresignedURL(ctx context.Context, bucketName string, objectName string, expiration time.Duration) (*entity.URL, error)
		UploadObject(ctx context.Context, bucketName string, objectName string, mimeType string, data []byte) (*entity.S3Key, error)
		DeleteObject(ctx context.Context, bucketName, objectName string) error
	}

	// Redis
	SessionRepository interface {
		AddSession(ctx context.Context, userID uint, accessToken string, expiration time.Duration) error
		DeleteSession(ctx context.Context, userID uint, accessToken string) error
		GetUserIDByAccessToken(ctx context.Context, accessToken string) (uint, error)
	}

	// Services

	// Auth service
	ServiceAuthRepository interface {
		CallLogin(ctx context.Context, email string, password string) (accessToken string, refreshToken string, err error)
		CallRefresh(ctx context.Context, refreshToken string) (accessToken string, err error)
		CallLogout(ctx context.Context, refreshToken string, accessToken string) error
		CallCreateUser(ctx context.Context, email string, username string, password string) (userID uint, accessToken string, refreshToken string, err error)
	}

	// Search service
	ServiceSearchRepository interface {
		CallSearchMedia(ctx context.Context, query string, limit, offset uint) ([]uint, error)
		CallSearchActors(ctx context.Context, query string, limit, offset uint) ([]uint, error)
	}

	// Payment gateway
	PaymentRepository interface {
		// Create payment and return payment ID
		CreatePayment(ctx context.Context, userID uint, amount string, description string) (string, error)

		// Approve payment by payment
		CapturePayment(ctx context.Context, payment *yoopayment.Payment) (*yoopayment.Payment, error)
	}
)
