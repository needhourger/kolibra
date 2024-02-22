package main

import (
	"kolibra/api"
	"kolibra/database"
	"log"
)


func main() {
	r := api.InitRouter()
	_,err := database.GetInstance()
	if err != nil {
		log.Panicf("Failed to connect to database: %v", err)
	}
	log.Printf("Database connected")
	// Listen and Server in 0.0.0.0:8080
	log.Printf("Server will running on 0.0.0.0:8080")
	r.Run(":8080")
}
