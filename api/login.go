package api

import (
	"kolibra/database"
	"kolibra/middleware"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	payload := LoginPayload{}
	c.BindJSON(&payload)
	log.Printf("Login Payload %v",payload)

	user, err := database.GetUserByUsername(payload.Username)
	if err != nil {
		c.JSON(404,gin.H{"message":"User not found"})
		return
	}

	if user.Password != payload.Password {
		c.JSON(403,gin.H{"message":"Invalid password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, middleware.JWTClaims{
		ID: user.ID,
		Username: user.Username,
		Role: user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		c.JSON(500,gin.H{"message":"Failed to generate token"})
		return
	}

	log.Printf("User %v logged in",user.Username)
	c.JSON(200,gin.H{"token":tokenString})
}