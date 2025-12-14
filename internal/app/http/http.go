package http

import (
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/adapter/aws"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/adapter/redis"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/adapter/yookassa"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/app"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/metrics"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/adapter/minio"
	db "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/adapter/postgres"
	grpc_auth "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/adapter/service/auth"
	grpc_search "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/adapter/service/search"
	cfg "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/config"
	srv "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http/handlers"
	rtr "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http/router"
	uc "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/usecase/http"
)

// Implements

// Minio
var _ uc.ObjectRepository = &minio.Minio{}
var _ uc.ObjectRepository = &aws.AWSS3{}

// Postgres
var _ uc.MediaRepository = &db.DataBase{}
var _ uc.UserRepository = &db.DataBase{}
var _ uc.ActorRepository = &db.DataBase{}
var _ uc.AssetRepository = &db.DataBase{}
var _ uc.AppealRepository = &db.DataBase{}
var _ uc.LikeRepository = &db.DataBase{}
var _ uc.GenreRepository = &db.DataBase{}

// Redis
var _ uc.SessionRepository = &redis.Redis{}

// Service

// Auth service
var _ uc.ServiceAuthRepository = &grpc_auth.AuthService{}

// Search service
var _ uc.ServiceSearchRepository = &grpc_search.SearchService{}

// Payment
var _ uc.PaymentRepository = &yookassa.Yookassa{}

func Run(logger logger.Logger, config cfg.Config) {

	// --- Connect to external services ---
	// Create Postgres connection
	dbURL := "postgres://" + config.APP_DB_USER + ":" + config.APP_DB_PASSWORD +
		"@" + config.POSTGRES_HOST + ":" + "5432" + "/" + config.POSTGRES_DB + "?sslmode=disable"
	var databaseAdapter app.DatabaseRepository = db.NewDataBase(logger, dbURL, config.DB_POOL_MAX_OPEN, config.DB_POOL_MAX_IDLE, config.DB_POOL_CONN_MAX_LIFETIME_MIN)
	err := databaseAdapter.Connect()
	if err != nil {
		logger.Fatal("Failed to connect to database: " + err.Error())
	}
	defer func() {
		if cerr := databaseAdapter.Close(); cerr != nil {
			logger.Error("Failed to close database adapter", logger.ToError(cerr))
		}
	}()

	// Create redis connection
	var cache app.Cache
	logger.Info("Connecting to redis", config.REDIS_HOST+":6379")
	redisClient := redis.NewRedis(logger, config.REDIS_HOST+":6379", "")
	defer func() {
		if cerr := redisClient.Close(); cerr != nil {
			logger.Error("Failed to close redis client", logger.ToError(cerr))
		}
	}()
	err = redisClient.CheckConnection()
	if err != nil {
		logger.Fatal("Failed to connect to Redis: " + err.Error())
	}
	cache = redisClient

	// Create s3 connection
	var s3 app.S3

	// AWS S3 required env variables:
	// - AWS_ACCESS_KEY_ID
	// - AWS_SECRET_ACCESS_KEY
	// - AWS_REGION
	// - AWS_S3_ENDPOINT
	logger.Info("Connecting to AWS S3")
	s3, err = aws.NewAWSS3(logger, config.AWS_S3_PUBLIC_URL)
	if err != nil {
		logger.Fatal("Failed to connect to AWS S3: " + err.Error())
	}
	logger.Info("AWS S3 connection established")

	// Connect Auth gRPC service
	var authService app.Service = grpc_auth.NewAuthService(logger, config.SERVICE_AUTH_SERVE_STRING)
	err = authService.Connect()
	if err != nil {
		logger.Fatal("Failed to connect to Auth gRPC service: " + err.Error())
	}
	defer authService.Close()

	// Connect Search gRPC service
	var searchService app.Service = grpc_search.NewSearchService(logger, config.SERVICE_SEARCH_SERVE_STRING)
	err = searchService.Connect()
	if err != nil {
		logger.Fatal("Failed to connect to Search gRPC service: " + err.Error())
	}
	defer searchService.Close()

	// --- Create repository level ---

	// Cast Postgres to MovieRepository
	mediaRepository, ok := databaseAdapter.(uc.MediaRepository)
	if !ok {
		logger.Fatal("Database can't be converted to MediaRepository")
	}
	// Cast Postgres to UserRepository
	userRepository, ok := databaseAdapter.(uc.UserRepository)
	if !ok {
		logger.Fatal("Database can't be converted to UserRepository")
	}

	// Cast Postgres to ActorRepository
	actorRepository, ok := databaseAdapter.(uc.ActorRepository)
	if !ok {
		logger.Fatal("Database can't be converted to ActorRepository")
	}
	// Cast Postgres to AssetRepository
	assetRepository, ok := databaseAdapter.(uc.AssetRepository)
	if !ok {
		logger.Fatal("Database can't be converted to AssetRepository")
	}

	// Cast Postgres to AppealRepository
	appealRepository, ok := databaseAdapter.(uc.AppealRepository)
	if !ok {
		logger.Fatal("Database can't be converted to AppealRepository")
	}

	// Cast Postgres to LikeRepository
	likeRepository, ok := databaseAdapter.(uc.LikeRepository)
	if !ok {
		logger.Fatal("Database can't be converted to LikeRepository")
	}

	// Cast Postgres to GenreRepository
	genreRepository, ok := databaseAdapter.(uc.GenreRepository)
	if !ok {
		logger.Fatal("Database can't be converted to GenreRepository")
	}

	// Cast Minio to ObjectRepository
	objectRepository, ok := s3.(uc.ObjectRepository)
	if !ok {
		logger.Fatal("Database can't be converted to ObjectRepository")
	}

	// Cast Redis to SessionRepository
	sessionRepository, ok := cache.(uc.SessionRepository)
	if !ok {
		logger.Fatal("Cache can't be converted to SessionRepository")
	}
	_ = sessionRepository // TODO: redis is unused!

	// Cast gRPC Auth service to AuthService
	authServiceRepository, ok := authService.(uc.ServiceAuthRepository)
	if !ok {
		logger.Fatal("gRPC Auth service can't be converted to AuthService")
	}

	// Cast gRPC Search service to SearchService
	searchServiceRepository, ok := searchService.(uc.ServiceSearchRepository)
	if !ok {
		logger.Fatal("gRPC Search service can't be converted to SearchService")
	}

	// Create Payment repository
	paymentRepository, err := yookassa.NewYookassa(
		logger,
		config.YOOKASSA_SHOP_ID,
		config.YOOKASSA_SECRET_KEY,
		config.YOOKASSA_RETURN_URL,
	)
	if err != nil {
		logger.Fatal("Failed to create Yookassa repository: " + err.Error())
	}

	// --- Create usecase level ---

	// Reusable usecases
	getObjectUseCase := uc.NewGetObjectUseCase(logger, objectRepository)
	getMediaUseCase := uc.NewGetMediaUseCase(logger, mediaRepository, getObjectUseCase, likeRepository)
	getUserUseCase := uc.NewGetUserMeUseCase(logger, userRepository, objectRepository)
	getActorUseCase := uc.NewGetActorUseCase(logger, actorRepository, getObjectUseCase)
	getGenreUseCase := uc.NewGetGenreUseCase(logger, genreRepository, mediaRepository)

	// Inject usecases into handler
	handler := handlers.NewHandlers(
		logger,
		uc.NewGetMediaRecommendationsUsecase(logger, mediaRepository, getMediaUseCase),
		getObjectUseCase,
		// Auth usecases
		uc.NewPostAuthSignInUsecase(logger, authServiceRepository),
		uc.NewGetAuthRefreshUseCase(logger, authServiceRepository),
		uc.NewPostAuthSignUpUsecase(logger, authServiceRepository),
		uc.NewGetAuthSignOutUsecase(logger, authServiceRepository),

		getUserUseCase,
		getActorUseCase,
		getMediaUseCase,
		uc.NewGetMediaWatchUseCase(logger, mediaRepository, getObjectUseCase, userRepository),
		uc.NewPostUserMeUpdateUseCase(logger, userRepository, getUserUseCase),
		uc.NewPostUserMeUpdateAvatarUseCase(logger, userRepository, objectRepository, assetRepository),
		uc.NewGetActorMediaUseCase(logger, actorRepository, getMediaUseCase),
		uc.NewGetMediaActorUseCase(logger, actorRepository, getActorUseCase),
		uc.NewPostUserMeUpdatePasswordUseCase(logger, userRepository),

		// Appeal usecases
		uc.NewGetAppealMyUseCase(logger, appealRepository),
		uc.NewPostAppealNewUseCase(logger, appealRepository),
		uc.NewGetAppealUseCase(logger, appealRepository),
		uc.NewPutAppealResolveUseCase(logger, appealRepository),
		// ----
		uc.NewPostAppealMessageUseCase(logger, appealRepository),
		uc.NewGetAppealMessageUseCase(logger, appealRepository),
		uc.NewGetAppealAllUseCase(logger, appealRepository),
		uc.NewGetSearchUseCase(logger, searchServiceRepository, getMediaUseCase, getActorUseCase),
		// Like usecases
		uc.NewGetMediaLikeUseCase(logger, likeRepository),
		uc.NewPutMediaLikeUseCase(logger, likeRepository),
		uc.NewDeleteMediaLikeUseCase(logger, likeRepository),
		// My media usecase
		uc.NewGetMediaMyUseCase(logger, mediaRepository, getMediaUseCase),
		// Genre usecases
		getGenreUseCase,
		uc.NewGetGenreAllUseCase(logger, genreRepository, getGenreUseCase),
		uc.NewGetGenreMediaUseCase(logger, mediaRepository, getMediaUseCase),
		// Episodes usecase
		uc.NewGetMediaEpisodesUseCase(logger, mediaRepository, getMediaUseCase),
		// Payment usecases
		uc.NewPostPaymentCompletedUsecase(logger, userRepository, paymentRepository),
		uc.NewPostPaymentNewUsecase(logger, paymentRepository, userRepository),
	)

	// Initialize JWT middleware engine
	common.InitJWT(config.SERVICE_HTTP_JWT_SECRET, config.SERVICE_HTTP_JWT_ACCESS_EXPIRATION, config.SERVICE_HTTP_JWT_REFRESH_EXPIRATION)
	logger.Info("JWT middleware initialized",
		logger.ToString("access_token_ttl", common.AccessTokenTTL.String()),
		logger.ToString("refresh_token_ttl", common.RefreshTokenTTL.String()),
	)
	// Inject handler into router
	router := rtr.InitRouter(handler, logger, config.SERVICE_HTTP_FRONTEND_URL)

	// Wrap with metrics middleware
	instrumentedRouter := metrics.HTTPMiddleware(metrics.ServiceHTTP, router)

	// Start HTTP server
	srv.StartServer(instrumentedRouter, config.SERVICE_HTTP_SERVESTRING, logger)
}
