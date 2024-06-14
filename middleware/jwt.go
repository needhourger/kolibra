package middleware

import (
	"errors"
	"kolibra/config"
	"kolibra/database"
	"log"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
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
	middleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "kolibra",
		Key:             []byte(config.Settings.Advance.JWTSecretKey),
		Timeout:         time.Duration(config.Settings.Advance.JWTTimeoutHours) * time.Hour,
		MaxRefresh:      time.Duration(config.Settings.Advance.JWTMaxRefreshHours) * time.Hour,
		IdentityKey:     config.Settings.Advance.JWTIdentityKey,
		TokenHeadName:   config.Settings.Advance.JWTHeadName,
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		SendCookie:      true,
		TimeFunc:        time.Now,
		PayloadFunc:     payloadFunc(),
		IdentityHandler: identityHandler(),
		Authenticator:   authenticator(),
		Authorizator:    authorizator(),
		Unauthorized:    unauthorizated(),
	})
	if err != nil {
		log.Fatalf("JWT init error %v", err)
	}
	return middleware
}

func identityHandler() func(c *gin.Context) any {
	return func(c *gin.Context) any {
		claims := jwt.ExtractClaims(c)
		log.Printf("JWT claims %v", claims[config.Settings.Advance.JWTIdentityKey])
		user, err := database.GetUserByID(claims[config.Settings.Advance.JWTIdentityKey].(string))
		if err != nil {
			log.Printf("JWT user fetch error: %v", err)
			return nil
		}
		return user
	}
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
			log.Println("JWT user convert failed")
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
		return user, nil
	}
}

func payloadFunc() func(data any) jwt.MapClaims {
	return func(data any) jwt.MapClaims {
		if v, ok := data.(*database.User); ok {
			return jwt.MapClaims{config.Settings.Advance.JWTIdentityKey: v.ID}
		}
		log.Println("JWT payloadFunc user convert failed")
		return jwt.MapClaims{}
	}
}
