package middleware

import (
	"kolibra/database"

	"github.com/gin-gonic/gin"
)

func AdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user")
		if user.(*database.User).Role != database.ADMIN {
			ctx.JSON(403, gin.H{"error": "Permission denied"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
