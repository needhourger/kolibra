package api

import (
	"kolibra/services/library"

	"github.com/gin-gonic/gin"
)

func ScanLibrary(c *gin.Context) {
	library.ScanLibrary()
	c.JSON(200, nil)
}
