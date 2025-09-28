package main

import (
	cfg "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/config"
	db "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/db"
	srv "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/server"
)

var config cfg.Config

func main() {
	// Load configuration
	config = cfg.Load()
	database := db.NewDataBase()
	server := srv.NewServer(&config, database)

	// Start the server
	server.Serve()
}
