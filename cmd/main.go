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
	redisClient := redis.NewRedis(logger, config.REDIS_HOST+":6379", "")
	defer redisClient.Close()
	err = redisClient.CheckConnection()
	if err != nil {
		logger.Fatal("Failed to connect to Redis: " + err.Error())
	}
	cache = redisClient

	// Create s3 connection
	var s3 uc.S3
	s3, err = minio.NewMinio(
		logger,
		config.MINIO_HOST,
		config.MINIO_ROOT_USER,
		config.MINIO_ROOT_PASSWORD,
		false,
	)
	if err != nil {
		logger.Fatal("Failed to connect to Minio: " + err.Error())
	}
	logger.Info("Minio connection established")

	// Usecases. Just follow openAPI order here

	// # Content

	// --- Get /movies ---
	movieRepository, ok := databaseAdapter.(uc.MovieRepository)
	if !ok {
		logger.Fatal("Database can't be converted to MovieRepository")
	}

	// --- Get /object ---
	objectRepository, ok := s3.(uc.ObjectRepository)
	if !ok {
		logger.Fatal("Database can't be converted to ObjectRepository")
	}

	// --- Get /auth/signin ---
	// Cast Postgres
	userRepository, ok := databaseAdapter.(uc.UserRepository)
	if !ok {
		logger.Fatal("Database can't be converted to UserRepository")
	}
	// Cast Redis
	sessionRepository, ok := cache.(uc.SessionRepository)
	if !ok {
		logger.Fatal("Cache can't be converted to SessionRepository")
	}
	tokenRepository, ok := databaseAdapter.(uc.TokenRepository)
	if !ok {
		logger.Fatal("Database can't be converted to TokenRepository")
	}

	// Inject usecases into handler
	handler := handlers.NewHandlers(
		logger,
		uc.NewGetMovieRecommendationsUsecase(logger, movieRepository),
		uc.NewGetObjectUsecase(logger, objectRepository),
		uc.NewPostAuthSignInUsecase(logger, userRepository, tokenRepository, sessionRepository),
		uc.NewGetAuthRefreshUseCase(logger, tokenRepository),
		uc.NewPostAuthSignUpUsecase(logger, userRepository, tokenRepository, sessionRepository),
	)

	// Initialize JWT middleware engine
	common.InitJWT(config.SERVER_JWT_SECRET, common.AccessTokenTTL, common.RefreshTokenTTL)
	// Inject handler into router
	router := srv.InitRouter(handler, logger, config.SERVER_FRONTEND_URL)
	srv.StartServer(router, config.SERVER_SERVE_STRING, logger)
}
