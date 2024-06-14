package api

import (
	"kolibra/services/library"

	"github.com/gin-gonic/gin"
)

func ScanLibrary(c *gin.Context) {
	go library.ScanLibrary()
	c.JSON(200, gin.H{"message": "Scan running background"})
}
