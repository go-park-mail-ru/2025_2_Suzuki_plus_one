package auth

import (
	db "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/adapter/postgres"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/app"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	cfg "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/config"
	srv "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/grpc/auth"
	uc "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/usecase/auth"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

var _ uc.TokenRepository = &db.DataBase{}
var _ uc.UserRepository = &db.DataBase{}

func Run() {

	// Load configuration
	config := cfg.Load()

	// --- Initialization ---

	// Initialize logger
	logger := logger.NewZapLogger(config.POPFILMS_ENVIRONMENT == "development")
	defer logger.Sync()

	logger.Info("Config loaded")

	// Create Postgres connection
	dbURL := "postgres://" + config.POSTGRES_USER + ":" + config.POSTGRES_PASSWORD +
		"@" + config.POSTGRES_HOST + ":" + "5432" + "/" + config.POSTGRES_DB + "?sslmode=disable"
	var databaseAdapter app.DatabaseRepository = db.NewDataBase(logger, dbURL)
	err := databaseAdapter.Connect()
	if err != nil {
		logger.Fatal("Failed to connect to database: " + err.Error())
	}
	defer databaseAdapter.Close()

	// Initialize JWT settings
	common.InitJWT(config.SERVER_JWT_SECRET, config.SERVER_JWT_ACCESS_EXPIRATION, config.SERVER_JWT_REFRESH_EXPIRATION)

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

	// Start gRPC server
	authHandler.StartGRPCServer(config.AUTH_SERVICE_SERVE_STRING)
}
