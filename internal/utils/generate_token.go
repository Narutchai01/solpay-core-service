package utils

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

func CreateJWTToken(secretKey string, claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

func GenerateAccesssToken(accountID uint) string {
	cfg := config.LoadConfig()
	secret := cfg.SECRET_JWT
	jwtExprise := cfg.JWT_EXPIRATION_HOURS

	hours, _ := strconv.Atoi(jwtExprise)
	cliams := jwt.MapClaims{
		"account_id": accountID,
		"exp":        time.Now().Add(time.Hour * time.Duration(hours)).Unix(),
	}
	token, err := CreateJWTToken(secret, cliams)
	if err != nil {
		slog.Error("GenerateAccesssToken: failed to create JWT token", "error", err)
		return ""
	}
	return token
}
