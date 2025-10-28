package main

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"

	db "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/adapter/inmemory"
	cfg "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/config"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller"
	rtr "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http/handlers"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http/middleware"
	uc "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/usecase"
)

func main() {

	// Load configuration
	config := cfg.Load()

	// Initialize logger
	logger := logger.NewZapLogger(config.POPFILMS_ENVIRONMENT == "development")
	defer logger.Sync()

	logger.Info("Config loaded")

	// TODO: clean arch
	// Create databaseAdapter connection
	databaseAdapter := db.NewDataBase()
	logger.Info("Database created")

	// Usecase
	usecase := uc.NewGetMoviesUsecase(databaseAdapter)

	// Inject usecase into router (pass as controller.GetMoviesUsecase interface)
	router := InitRouter(usecase, logger)

	// Create server
	// TODO: move to server package
	// TODO: rollback to zap logger without interface, because I cannot see debug logs properly
	startServer(router, logger, databaseAdapter, &config)
}

func startServer(router http.Handler, l logger.Logger, database *db.DataBase, config *cfg.Config) {
	// Create server
	server := http.NewServeMux()
	server.Handle("/", router)

	// Start the server
	if err := http.ListenAndServe(config.SERVER_SERVE_STRING, server); err != nil {
		l.Fatal("Server not started 'cause of error", l.ToError(err))
	}
	l.Info("Server started serving")
}

func InitRouter(movieUC controller.GetMoviesUsecase, l logger.Logger) http.Handler {
	r := rtr.NewRouter("", l)

	// Custom middleware
	r.Use(middleware.GetLogging(l))

	// Handlers
	h := handlers.NewHandlers(movieUC, l)

	// Register routes here
	r.Handle(rtr.GET, "/movies", h.GetMovies)

	return r
}
