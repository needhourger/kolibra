package main

import (
	"fmt"
	"kolibra/api"
	"kolibra/config"
	"kolibra/database"
	"kolibra/services/reader"
	"log"
)

func main() {
	// load config
	config.LoadConfig()
	reader.CreateReaderCache()

	// Connect to database
	err := database.InitDatabase()
	if err != nil {
		log.Panicf("Failed to connect to database: %v", err)
	}
	log.Printf("Database connected")

	// Set up router
	r := api.InitRouter()
	address := fmt.Sprintf("%s:%d", config.Config.Host, config.Config.Port)
	log.Printf("Server will running on %s", address)
	r.Run(address)
}
