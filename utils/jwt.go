package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "your_secret_key"

func GenerateToken(email string, userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 2).Unix(), // Token expires in 2 hours
	})

	return token.SignedString([]byte(secretKey)) // Sign the token with the secret key
}

func VerifyToken(token string) (error, int64) {
	parasedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		// 	return nil, jwt.ErrSignatureInvalid
		// }

		// commented block is the same as below
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("Unexpected signing method")
		}

		return []byte(secretKey), nil
	})
	if err != nil {
		return errors.New("Invalid token: " + err.Error()), 0
	}

	if !parasedToken.Valid {
		return errors.New("Invalid token"), 0
	}

	claims, ok := parasedToken.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("Invalid token claims"), 0
	}

	userId := int64(claims["userId"].(float64))

	return nil, userId
}
