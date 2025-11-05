package main

import (
	"log"

	"github.com/RadekKusiak71/places-app/config"
	"github.com/RadekKusiak71/places-app/database"
	"github.com/RadekKusiak71/places-app/server"
)

func main() {
	config.InitConfig()

	apiServer := server.NewAPIServer(config.Config.PORT, database.New())

	if err := apiServer.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
