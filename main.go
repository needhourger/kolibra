package main

import (
	"fmt"
	"kolibra/api"
	"kolibra/config"
	DB "kolibra/models"
	"kolibra/services/reader"
	"log"
)

func main() {
	// load config
	config.LoadConfig()
	reader.CreateReaderCache()

	// Connect to database
	err := DB.InitDatabase()
	if err != nil {
		log.Panicf("Failed to connect to database: %v", err)
	}
	log.Printf("Database connected")

	// Set up router
	engine := api.InitGinEngine()
	address := fmt.Sprintf("%s:%d", config.Settings.Host, config.Settings.Port)
	log.Printf("Server will running on %s", address)
	engine.Run(address)
}
