package middleware

import (
	"kolibra/database"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var SkipAuth = []string{"/api/login","/api/ping","/api/sign"}

type JWTClaims struct {
	ID uint `json:"id"`
	Username string `json:"username"`
	Role string `json:"role"`
	jwt.StandardClaims
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _,path := range SkipAuth {
			if path == ctx.FullPath() {
				ctx.Next()
				return
			}
		}
		// Get the token from the header
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(401, gin.H{"error": "Authorization header is required"})
			ctx.Abort()
			return
		}

		// Parse the token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil {
			ctx.JSON(401, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		// Check if the token is valid
		claims, ok := token.Claims.(*JWTClaims)
		if !ok || !token.Valid {
			ctx.JSON(401, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		user, err := database.GetUserByID(claims.ID)
		if err != nil {
			ctx.JSON(401, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		// Set the user ID in the context
		ctx.Set("user", user)
		ctx.Next()
	}
}