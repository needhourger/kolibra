//go:build production

package config

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func InitLogFormat() {
	file, err := os.OpenFile("data/kolibra.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()
	gin.SetMode(gin.ReleaseMode)

	log.SetOutput(file)
}
