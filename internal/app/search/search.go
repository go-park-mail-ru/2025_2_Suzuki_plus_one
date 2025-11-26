package search

import (
	db "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/adapter/postgres"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/app"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	cfg "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/config"
	srv "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/grpc/search"
	uc "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/usecase/search"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/metrics"
)

var _ uc.MediaRepository = &db.DataBase{}
var _ uc.ActorRepository = &db.DataBase{}

func Run(logger logger.Logger, config cfg.Config) {

	// --- Initialization ---

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
	common.InitJWT(
		config.SERVICE_HTTP_JWT_SECRET,
		config.SERVICE_HTTP_JWT_ACCESS_EXPIRATION,
		config.SERVICE_HTTP_JWT_REFRESH_EXPIRATION,
	)

	// --- Create repository level ---

	// Cast Postgres to MediaRepository
	mediaRepo, ok := databaseAdapter.(uc.MediaRepository)
	if !ok {
		logger.Fatal("Database can't be converted to MediaRepository")
	}

	// Cast Postgres to ActorRepository
	actorRepo, ok := databaseAdapter.(uc.ActorRepository)
	if !ok {
		logger.Fatal("Database can't be converted to ActorRepository")
	}

	// --- Create usecase level ---
	searchActorUsecase := uc.NewSearchActorUsecase(logger, actorRepo)
	searchMediaUsecase := uc.NewSearchMediaUsecase(logger, mediaRepo)

	// --- Create delivery level ---
	searchHandler := srv.NewSearchServer(logger, searchMediaUsecase, searchActorUsecase)

	// Add metrics middleware
	searchHandler.Middleware = append(searchHandler.Middleware,
		metrics.GRPCServerMetricsInterceptor(metrics.ServiceSearch),
	)

	// Start gRPC server
	searchHandler.StartGRPCServer(config.SERVICE_SEARCH_SERVE_STRING)
}
