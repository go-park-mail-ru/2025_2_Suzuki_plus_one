package main

import (
	db "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/db"
	srv "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/server"
)

var (
	address  = ":8080"
	database = db.NewDataBase()
	server   = srv.NewServer(address, database)
)

func main() {
	// Start the server
	server.Serve()
}
