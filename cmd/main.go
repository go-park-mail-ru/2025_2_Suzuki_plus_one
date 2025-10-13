package main

import (
	"log"

	cfg "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/config"
	db "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/db"
	srv "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/server"
)

func main() {
	// Load configuration
	config := cfg.Load()
	log.Println("Config loaded")

	database := db.NewDataBase()
	log.Println("Database created")

	server := srv.NewServer(&config, database)
	log.Println("Server created")

	// Start the server
	server.Serve()
	log.Println("Server started serving")
}
