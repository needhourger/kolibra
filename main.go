package main

import (
	"fmt"
	"kolibra/api"
	"kolibra/config"
	"kolibra/database/dao"
	"kolibra/database/model"
	"kolibra/services/reader"
	"log"
)

func main() {
	// load config
	config.LoadConfig()
	reader.CreateReaderCache()

	// Connect to database
	model.InitDatabase()
	dao.InitDAO()
	log.Printf("Database connected")

	// Set up router
	engine := api.InitGinEngine()
	address := fmt.Sprintf("%s:%d", config.Settings.Host, config.Settings.Port)
	log.Printf("Server will running on %s", address)
	engine.Run(address)
}
