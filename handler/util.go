package handler

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func sendError(c *gin.Context, status int, err error) {
	log.Println(err)
	c.JSON(status, gin.H{
		"message": err.Error(),
	})
}

type claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func generateToken(userID uint) (string, error) {
	payload := claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			Issuer:    "course-api",
		},
	}
	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return claim.SignedString([]byte(secretKey))
}
