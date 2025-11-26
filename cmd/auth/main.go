package main

import (
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/app/auth"
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
	logger.Info("Auth service: Config loaded")

	// Run HTTP metrics service
	go metrics.Serve(logger, config.SERVICE_AUTH_METRICS_SERVE_STRING)
	logger.Info("Auth service: Metrics service started")


	// Run gRPC auth service
	auth.Run(logger, config)
}
