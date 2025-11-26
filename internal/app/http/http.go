package http

import (
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/adapter/redis"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/app"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"

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

func Run(logger logger.Logger, config cfg.Config) {

	// --- Connect to external services ---
	// Create Postgres connection
	dbURL := "postgres://" + config.POSTGRES_USER + ":" + config.POSTGRES_PASSWORD +
		"@" + config.POSTGRES_HOST + ":" + "5432" + "/" + config.POSTGRES_DB + "?sslmode=disable"
	var databaseAdapter app.DatabaseRepository = db.NewDataBase(logger, dbURL)
	err := databaseAdapter.Connect()
	if err != nil {
		logger.Fatal("Failed to connect to database: " + err.Error())
	}
	defer databaseAdapter.Close()

	// Create redis connection
	var cache app.Cache
	logger.Info("Connecting to redis", config.REDIS_HOST+":6379")
	redisClient := redis.NewRedis(logger, config.REDIS_HOST+":6379", "")
	defer redisClient.Close()
	err = redisClient.CheckConnection()
	if err != nil {
		logger.Fatal("Failed to connect to Redis: " + err.Error())
	}
	cache = redisClient

	// Create s3 connection
	var s3 app.S3

	// URL for media files will be like http(s)://POPFILMS_SERVICE_HTTP_FRONTEND_URL/bucketName/objectName
	minioServePrefix := config.SERVICE_HTTP_FRONTEND_URL
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

	// Connect Auth gRPC service
	var authService app.Service

	authService = grpc_auth.NewAuthService(logger, config.SERVICE_AUTH_SERVE_STRING)
	err = authService.Connect()
	if err != nil {
		logger.Fatal("Failed to connect to Auth gRPC service: " + err.Error())
	}
	defer authService.Close()

	// Connect Search gRPC service
	var searchService app.Service

	searchService = grpc_search.NewSearchService(logger, config.SERVICE_SEARCH_SERVE_STRING)
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

	// --- Create usecase level ---

	// Reusable usecases
	getObjectUseCase := uc.NewGetObjectUseCase(logger, objectRepository)
	getMediaUseCase := uc.NewGetMediaUseCase(logger, mediaRepository, getObjectUseCase)
	getUserUseCase := uc.NewGetUserMeUseCase(logger, userRepository, objectRepository)
	getActorUseCase := uc.NewGetActorUseCase(logger, actorRepository, getObjectUseCase)
	getGenreUseCase := uc.NewGetGenreUseCase(logger, genreRepository, mediaRepository, getMediaUseCase)

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
		uc.NewGetMediaWatchUseCase(logger, mediaRepository, getObjectUseCase),
		uc.NewPostUserMeUpdateUseCase(logger, userRepository, getUserUseCase),
		uc.NewPostUserMeUpdateAvatarUseCase(logger, userRepository, objectRepository, assetRepository),
		uc.NewGetActorMediaUseCase(logger, actorRepository, getMediaUseCase),
		uc.NewGetMediaActorUseCase(logger, actorRepository, getActorUseCase),
		uc.NewPostUserMeUpdatePasswordUseCase(logger, userRepository),

		// Appeal usecases
		uc.NewGetAppealMyUseCase(logger, appealRepository),
		uc.NewPostAppealNewUseCase(logger, appealRepository),    //done
		uc.NewGetAppealUseCase(logger, appealRepository),        //done
		uc.NewPutAppealResolveUseCase(logger, appealRepository), //done
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
	)

	// Initialize JWT middleware engine
	common.InitJWT(config.SERVICE_HTTP_JWT_SECRET, config.SERVICE_HTTP_JWT_ACCESS_EXPIRATION, config.SERVICE_HTTP_JWT_REFRESH_EXPIRATION)
	logger.Info("JWT middleware initialized",
		logger.ToString("access_token_ttl", common.AccessTokenTTL.String()),
		logger.ToString("refresh_token_ttl", common.RefreshTokenTTL.String()),
	)
	// Inject handler into router
	router := rtr.InitRouter(handler, logger, config.SERVICE_HTTP_FRONTEND_URL)
	srv.StartServer(router, config.SERVICE_HTTP_SERVESTRING, logger)
}
