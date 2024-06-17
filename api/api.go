package api

import (
	"kolibra/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func initRouter(engine *gin.Engine, jwtMiddleware *jwt.GinJWTMiddleware) {
	// apiBase := r.Group("/apiBase", middleware.JWTAuthMiddleware())
	apiBase := engine.Group("/api")

	// No auth api
	apiBase.POST("/auth", jwtMiddleware.LoginHandler)
	apiBase.POST("/sign", Sign)

	// Auth api
	// Book api
	bookApi := apiBase.Group("/books")
	// bookApi.Use(jwtMiddleware.MiddlewareFunc())
	bookApi.GET("/", GetAllBooks)
	bookApi.GET("/:id", GetBook)
	bookApi.DELETE("/:id", DeleteBookByID)
	bookApi.GET("/:id/chapters", GetBookChapters)
	bookApi.GET("/:id/chapters/:cid", GetBookChapter)
	bookApi.GET("/:id/chapters/:cid/content", GetChapterContent)
	// Library api
	libraryApi := apiBase.Group("/library")
	// libraryApi.Use(jwtMiddleware.MiddlewareFunc())
	libraryApi.GET("/scan", ScanLibrary)
}

func InitGinEngine() *gin.Engine {
	engine := gin.Default()

	jwtMiddleware := middleware.InitJWTMiddleware()
	initRouter(engine, jwtMiddleware)

	return engine
}
