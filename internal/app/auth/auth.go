package auth

import (
	db "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/adapter/postgres"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/app"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	cfg "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/config"
	srv "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/grpc/auth"
	uc "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/usecase/auth"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/metrics"
)

var _ uc.TokenRepository = &db.DataBase{}
var _ uc.UserRepository = &db.DataBase{}

func Run(logger logger.Logger, config cfg.Config) {

	// --- Initialization ---

	// Create Postgres connection
	dbURL := "postgres://" + config.APP_DB_USER + ":" + config.APP_DB_PASSWORD +
		"@" + config.POSTGRES_HOST + ":" + "5432" + "/" + config.POSTGRES_DB + "?sslmode=disable"
	var databaseAdapter app.DatabaseRepository = db.NewDataBase(logger, dbURL, config.DB_POOL_MAX_OPEN, config.DB_POOL_MAX_IDLE, config.DB_POOL_CONN_MAX_LIFETIME_MIN)
	err := databaseAdapter.Connect()
	if err != nil {
		logger.Fatal("Failed to connect to database: " + err.Error())
	}
	defer databaseAdapter.Close()

	// Initialize JWT settings
	common.InitJWT(config.SERVICE_HTTP_JWT_SECRET, config.SERVICE_HTTP_JWT_ACCESS_EXPIRATION, config.SERVICE_HTTP_JWT_REFRESH_EXPIRATION)

	// --- Create repository level ---

	// Cast Postgres to UserRepository
	userRepo, ok := databaseAdapter.(uc.UserRepository)
	if !ok {
		logger.Fatal("Database can't be converted to UserRepository")
	}

	// Cast Postgres to TokenRepository
	tokenRepo, ok := databaseAdapter.(uc.TokenRepository)
	if !ok {
		logger.Fatal("Database can't be converted to TokenRepository")
	}

	// --- Create usecase level ---
	loginUsecase := uc.NewLoginUsecase(logger, userRepo, tokenRepo)
	refreshUsecase := uc.NewRefreshUsecase(logger, userRepo, tokenRepo)
	logoutUsecase := uc.NewLogoutUsecase(logger, userRepo, tokenRepo)
	createUserUsecase := uc.NewCreateUserUsecase(logger, userRepo, tokenRepo)

	// --- Create delivery level ---
	authHandler := srv.NewAuthServer(logger, loginUsecase, refreshUsecase, logoutUsecase, createUserUsecase)

	// Add metrics middleware
	authHandler.Middleware = append(authHandler.Middleware,
		metrics.GRPCServerMetricsInterceptor(metrics.ServiceAuth),
	)
	// Start gRPC server
	authHandler.StartGRPCServer(config.SERVICE_AUTH_SERVE_STRING)
}
