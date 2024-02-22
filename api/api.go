package api

import (
	"github.com/gin-gonic/gin"
)

var skipAuth = []string{"/api/login","/api/ping","/api/sign"}

func InitRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api",JWTAuthMiddleware())
	api.GET("/ping", Ping)

	api.POST("/login", Login)
	api.POST("/sign", Sign)

	return r
}