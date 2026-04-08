package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(adminID string, username string, secret string) (string, error) {
	claims := jwt.MapClaims{
		"sub":      adminID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
