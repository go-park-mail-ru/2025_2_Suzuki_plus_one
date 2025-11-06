package main

import (
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/adapter/redis"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/adapter/minio"
	db "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/adapter/postgres"
	cfg "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/config"
	srv "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http/handlers"
	uc "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/usecase"
)

// Implements

// Minio
var _ uc.ObjectRepository = &minio.Minio{}

// Postgres
var _ uc.MediaRepository = &db.DataBase{}
var _ uc.UserRepository = &db.DataBase{}
var _ uc.TokenRepository = &db.DataBase{}
var _ uc.ActorRepository = &db.DataBase{}

// Redis
var _ uc.SessionRepository = &redis.Redis{}

func main() {

	// Load configuration
	config := cfg.Load()

	// Initialize logger
	logger := logger.NewZapLogger(config.POPFILMS_ENVIRONMENT == "development")
	defer logger.Sync()

	logger.Info("Config loaded")

	// Create Postgres connection
	dbURL := "postgres://" + config.POSTGRES_USER + ":" + config.POSTGRES_PASSWORD +
		"@" + config.POSTGRES_HOST + ":" + "5432" + "/" + config.POSTGRES_DB + "?sslmode=disable"
	var databaseAdapter uc.Repository = db.NewDataBase(logger, dbURL)
	err := databaseAdapter.Connect()
	if err != nil {
		logger.Fatal("Failed to connect to database: " + err.Error())
	}
	defer databaseAdapter.Close()

	// Create redis connection
	var cache uc.Cache
	logger.Info("Connecting to redis", config.REDIS_HOST+":6379")
	redisClient := redis.NewRedis(logger, config.REDIS_HOST+":6379", "")
	defer redisClient.Close()
	err = redisClient.CheckConnection()
	if err != nil {
		logger.Fatal("Failed to connect to Redis: " + err.Error())
	}
	cache = redisClient

	// Create s3 connection
	var s3 uc.S3

	// URL for media files will be like http(s)://SERVER_FRONTEND_URL/bucketName/objectName
	minioServePrefix := config.SERVER_FRONTEND_URL
	// must not end with /
	if minioServePrefix[len(minioServePrefix)-1] == '/' {
		minioServePrefix = minioServePrefix[:len(minioServePrefix)-1]
	}
	s3, err = minio.NewMinio(
		logger,
		config.MINIO_INTERNAL_HOST,
		config.MINIO_EXTERNAL_HOST,
		config.MINIO_ROOT_USER,
		config.MINIO_ROOT_PASSWORD,
		false,
	)
	if err != nil {
		logger.Fatal("Failed to connect to Minio: " + err.Error())
	}
	logger.Info("Minio connection established")

	// Cast Postgres to MovieRepository
	movieRepository, ok := databaseAdapter.(uc.MediaRepository)
	if !ok {
		logger.Fatal("Database can't be converted to MediaRepository")
	}
	// Cast Postgres to UserRepository
	userRepository, ok := databaseAdapter.(uc.UserRepository)
	if !ok {
		logger.Fatal("Database can't be converted to UserRepository")
	}
	// Cast Postgres to TokenRepository
	tokenRepository, ok := databaseAdapter.(uc.TokenRepository)
	if !ok {
		logger.Fatal("Database can't be converted to TokenRepository")
	}
	// Cast Postgres to ActorRepository
	actorRepository, ok := databaseAdapter.(uc.ActorRepository)
	if !ok {
		logger.Fatal("Database can't be converted to ActorRepository")
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

	// Reusable usecases
	getObjectUseCase := uc.NewGetObjectUseCase(logger, objectRepository)
	getMediaUseCase := uc.NewGetMediaUseCase(logger, movieRepository, actorRepository, getObjectUseCase)
	getUserUseCase := uc.NewGetUserMeUseCase(logger, userRepository, sessionRepository, objectRepository)

	// Inject usecases into handler
	handler := handlers.NewHandlers(
		logger,
		uc.NewGetMovieRecommendationsUsecase(logger, movieRepository, getMediaUseCase),
		getObjectUseCase,
		uc.NewPostAuthSignInUsecase(logger, userRepository, tokenRepository, sessionRepository),
		uc.NewGetAuthRefreshUseCase(logger, tokenRepository),
		uc.NewPostAuthSignUpUsecase(logger, userRepository, tokenRepository, sessionRepository),
		uc.NewGetAuthSignOutUsecase(logger, tokenRepository, sessionRepository),
		getUserUseCase,
		uc.NewGetActorUseCase(logger, actorRepository, getMediaUseCase, getObjectUseCase),
		getMediaUseCase,
		uc.NewGetMediaWatchUseCase(logger, movieRepository, getObjectUseCase),
		uc.NewPostUserMeUpdateUseCase(logger, userRepository, getUserUseCase),
	)

	// Initialize JWT middleware engine
	common.InitJWT(config.SERVER_JWT_SECRET, config.SERVER_JWT_ACCESS_EXPIRATION, config.SERVER_JWT_REFRESH_EXPIRATION)
	logger.Info("JWT middleware initialized",
		logger.ToString("access_token_ttl", common.AccessTokenTTL.String()),
		logger.ToString("refresh_token_ttl", common.RefreshTokenTTL.String()),
	)
	// Inject handler into router
	router := srv.InitRouter(handler, logger, config.SERVER_FRONTEND_URL)
	srv.StartServer(router, config.SERVER_SERVE_STRING, logger)
}
