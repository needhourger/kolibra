package api

import (
	"kolibra/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// api := r.Group("/api", middleware.JWTAuthMiddleware())
	api := r.Group("/api")
	// Test api
	api.GET("/ping", Ping)
	api.GET("/pings", middleware.AdminAuth(), Ping)

	// No auth api
	api.POST("/login", Login)
	api.POST("/sign", Sign)

	// Auth api
	api.GET("/books", GetAllBooks)
	api.GET("/books/:id", GetBook)
	api.GET("/books/:id/chapters", GetBookChapters)
	api.GET("/library/scan", ScanLibrary)

	return r
}
