package middleware

import (
	"errors"
	"kolibra/config"
	"kolibra/database"
	"log"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	ErrUserNotFound    = errors.New("User not found")
	ErrPasswordInvalid = errors.New("Password invalid")
)

func InitJWTMiddleware() *jwt.GinJWTMiddleware {
	middleware := &jwt.GinJWTMiddleware{
		Realm:         "Kolibra",
		Key:           []byte(config.Settings.Advance.JWTSecretKey),
		Timeout:       time.Duration(config.Settings.Advance.JWTTimeoutHours) * time.Hour,
		MaxRefresh:    time.Duration(config.Settings.Advance.JWTMaxRefreshHours) * time.Hour,
		TokenHeadName: config.Settings.Advance.JWTHeadName,
		SendCookie:    true,
		PayloadFunc:   payloadFunc(),
		Authenticator: authenticator(),
		Authorizator:  authorizator(),
		Unauthorized:  unauthorizated(),
		TokenLookup:   "header:Authorization, query: token, cookie: jwt",
		TimeFunc:      time.Now,
	}
	return middleware
}

// Api permission check callback
func unauthorizated() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{"message": message})
	}
}

// Api permission check function
func authorizator() func(data any, c *gin.Context) bool {
	return func(data any, c *gin.Context) bool {
		user, ok := data.(*database.User)
		if !ok {
			return false
		}
		requestPath := c.Request.URL.Path
		if strings.Contains(requestPath, "/api/admin") && user.Role != database.ADMIN {
			return false
		}
		return true
	}
}

// Auth api handler
func authenticator() func(c *gin.Context) (any, error) {
	return func(c *gin.Context) (any, error) {
		var payload LoginPayload
		if err := c.Bind(&payload); err != nil {
			return "", jwt.ErrMissingLoginValues
		}
		user, err := database.GetUserByUsername(payload.Username)
		if err != nil {
			return "", ErrUserNotFound
		}
		if user.Password != payload.Password {
			return "", ErrPasswordInvalid
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
