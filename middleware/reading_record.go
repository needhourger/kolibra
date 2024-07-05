package middleware

import (
	DB "kolibra/database"
	"log"

	"github.com/gin-gonic/gin"
)

func RecordReading() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := GetUserFromJWT(ctx)
		bookID := ctx.Param("id")
		chapterID := ctx.Param("cid")
		record := &DB.ReadingRecord{
			UserID:    user.ID,
			BookID:    bookID,
			ChapterID: chapterID,
		}
		log.Printf("Created record: %v", record)
	}
}
