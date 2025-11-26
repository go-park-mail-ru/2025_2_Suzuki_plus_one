package main

import (
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/app/search"
	cfg "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/config"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/metrics"
)

func main() {
	// Load configuration
	config := cfg.Load()

	// Initialize logger
	logger := logger.NewZapLogger(config.ENVIRONMENT == "development")
	defer logger.Sync()
	logger.Info("Search service: Config loaded")

	// Run HTTP metrics service
	go metrics.Serve(logger, config.SERVICE_SEARCH_METRICS_SERVE_STRING)
	logger.Info("Search service: Metrics service started")

	// Run gRPC auth service
	search.Run(logger, config)
}
