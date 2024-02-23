package api

import (
	"kolibra/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api", middleware.JWTAuthMiddleware())
	api.GET("/ping", Ping)

	api.POST("/login", Login)
	api.POST("/sign", Sign)

	return r
}
