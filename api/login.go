package api

import (
	"kolibra/database"
	"log"

	"github.com/gin-gonic/gin"
)

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	payload := LoginPayload{}
	c.BindJSON(&payload)
	log.Printf("Login Payload %v", payload)

	user, err := database.GetUserByUsername(payload.Username)
	if err != nil {
		c.JSON(404, gin.H{"message": "User not found"})
		return
	}

	if user.Password != payload.Password {
		c.JSON(403, gin.H{"message": "Invalid password"})
		return
	}

	//todo: golang jwt auth

	log.Printf("User %v logged in", user.Username)
	c.JSON(200, gin.H{"token": ""})
}
