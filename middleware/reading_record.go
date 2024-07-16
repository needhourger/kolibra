package middleware

import (
	"kolibra/database/model"
	"log"

	"github.com/gin-gonic/gin"
)

func RecordReading() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := GetUserFromContext(ctx)
		bookID := ctx.Param("id")
		chapterID := ctx.Param("cid")
		record := &model.ReadingRecord{
			UserID:    user.ID,
			BookID:    bookID,
			ChapterID: chapterID,
		}
		model.CreateOrUpdateReadingRecord(record)
		log.Printf("Created record: %v", record)
	}
}
