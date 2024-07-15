package api

import (
	"kolibra/middleware"
	embedFs "kolibra/static"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func initRouter(engine *gin.Engine, jwtMiddleware *jwt.GinJWTMiddleware) {
	apiBase := engine.Group("/api")

	// No auth api
	apiBase.POST("/auth", jwtMiddleware.LoginHandler)
	apiBase.POST("/sign", Sign)

	// Auth api
	// Book api
	bookApi := apiBase.Group("/books")
	bookApi.Use(jwtMiddleware.MiddlewareFunc())
	bookApi.GET("/", GetAllBooks)
	bookApi.GET("/:id", GetBook)
	bookApi.DELETE("/:id", DeleteBookByID)
	bookApi.POST("/:id", UpdateBook)
	bookApi.GET("/:id/chapters", GetBookChapters)
	bookApi.GET("/:id/chapters/:cid", GetBookChapter)
	bookApi.GET("/:id/chapters/:cid/content", middleware.RecordReading(), GetChapterContent)

	// Reading record api
	readingRecordApi := apiBase.Group("/reading_record")
	readingRecordApi.Use(jwtMiddleware.MiddlewareFunc())
	readingRecordApi.GET("/:bookID", GetBookReadingRecord)
	// Library api
	libraryApi := apiBase.Group("/library")
	libraryApi.Use(jwtMiddleware.MiddlewareFunc())
	libraryApi.GET("/scan", ScanLibrary)
}

func initStatics(engine *gin.Engine) {
	engine.Use(static.Serve("/", static.EmbedFolder(embedFs.Embed, "dist")))
	engine.NoRoute(func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	})
}

func InitGinEngine() *gin.Engine {
	engine := gin.Default()

	jwtMiddleware := middleware.InitJWTMiddleware()
	initRouter(engine, jwtMiddleware)
	initStatics(engine)

	return engine
}
