package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken() {
	jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"email": "",
		"userId":"",
		"exp": time.Now().Add(time.Hour *2).Unix()
	})
}