package api

import (
	"kolibra/database"
	"log"

	"github.com/gin-gonic/gin"
)

type SignPayload struct {
	Email string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func Sign(c *gin.Context) {
	payload := SignPayload{}
	c.BindJSON(&payload)
	log.Printf("Sign Payload %v",payload)

	exist := database.CheckUserByName(payload.Username)
	if exist {
		c.JSON(403,gin.H{"message":"Username already exists"})
		return
	}

	user := database.User{
		Username: payload.Username,
		Password: payload.Password,
		Email: payload.Email,
		Role: "user",
	}

	err := database.CreateUser(&user)
	if err != nil {
		c.JSON(500,gin.H{"message":"Failed to create user"})
		return
	}
	c.JSON(200,gin.H{"message":"User created successfully"})
}