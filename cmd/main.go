package main

import (
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"

	db "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/adapter/inmemory"
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

	// Create databaseAdapter connection
	databaseAdapter := db.NewDataBase(logger)
	logger.Info("Database created")

	// Usecase
	getMovies := uc.NewGetMoviesUsecase(databaseAdapter)

	// Inject usecase into handler
	handler := handlers.NewHandlers(getMovies, logger)

	// Inject handler into router
	router := srv.InitRouter(handler, logger, config.SERVER_FRONTEND_URL)

	srv.StartServer(router, config.SERVER_SERVE_STRING, logger)
}
