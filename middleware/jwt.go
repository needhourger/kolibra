package middleware

import (
	"kolibra/config"
	"kolibra/database"
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func InitJWTMiddleware() *jwt.GinJWTMiddleware {
	middleware := &jwt.GinJWTMiddleware{
		Realm:         "Kolibra",
		Key:           []byte(config.Settings.Advance.JWTSecretKey),
		Timeout:       time.Duration(config.Settings.Advance.JWTTimeoutHours) * time.Hour,
		MaxRefresh:    time.Duration(config.Settings.Advance.JWTMaxRefreshHours) * time.Hour,
		TokenHeadName: config.Settings.Advance.JWTHeadName,
		PayloadFunc:   payloadFunc(),
		Authenticator: authenticator(),
		Authorizator:  authorizator(),
		Unauthorized:  unauthorizated(),
		TokenLookup:   "header:Authorization, query: token, cookie: jwt",
		TimeFunc:      time.Now,
	}
	return middleware
}

func unauthorizated() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{"message": message})
	}
}

func authorizator() func(data any, c *gin.Context) bool {
	return func(data any, c *gin.Context) bool {
		if v, ok := data.(database.User); ok && v.Role == database.ADMIN {
			return true
		}
		return false
	}
}

func authenticator() func(c *gin.Context) (any, error) {
	return func(c *gin.Context) (any, error) {
		var payload LoginPayload
		if err := c.Bind(&payload); err != nil {
			return "", jwt.ErrMissingLoginValues
		}
		user, err := database.GetUserByUsername(payload.Username)
		if err != nil {
			c.JSON(404, gin.H{"message": "User not found"})
		}
		if user.Password != payload.Password {
			c.JSON(403, gin.H{"message": "Invalid password"})
		}
		log.Printf("User %v logged in", user)
		return &user, nil
	}
}

func payloadFunc() func(data any) jwt.MapClaims {
	return func(data any) jwt.MapClaims {
		if v, ok := data.(*database.User); ok {
			return jwt.MapClaims{"user": v}
		}
		return jwt.MapClaims{}
	}
}
