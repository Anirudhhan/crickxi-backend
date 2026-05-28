package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/nyaruka/phonenumbers"
	"golang.org/x/crypto/bcrypt"
)

var accessSecret = []byte(os.Getenv("ACCESS_SECRET"))

func ValidatePhoneNumber(phone string) (string, error) {
	num, err := phonenumbers.Parse(phone, "IN")
	if err != nil {
		return "", err
	}
	if !phonenumbers.IsValidNumber(num) {
		return "", errors.New("invalid phone number")
	}

	return phonenumbers.Format(num, phonenumbers.E164), nil
}

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

func GenerateAccessToken(userID string, playerStatsID string, sessionID string) (string, error) {
	claims := jwt.MapClaims{
		"uid": userID,
		"pid": playerStatsID,
		"sid": sessionID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessSecret)
}

func ValidateAccessToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return accessSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}
