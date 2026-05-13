package utils

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var accessSecret = []byte(os.Getenv("ACCESS_SECRET"))

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	return string(hashed), err
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}

func ErrorResponse(ctx *gin.Context, status int, err error, message string) {
	fmt.Println("error: ", err.Error())
	ctx.JSON(status, gin.H{
		"error": message,
	})
}
