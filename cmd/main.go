package main

import (
	"log"

	"go.uber.org/zap"

	cfg "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/config"
	db "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/db"
	srv "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/server"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	// Load configuration
	config := cfg.Load()
	logger.Info("Config loaded")

	database := db.NewDataBase()
	logger.Info("Database created")

	server := srv.NewServer(&config, database, logger)
	logger.Info("Server created")

	// Start the server
	server.Serve()
	logger.Info("Server started serving")
}
