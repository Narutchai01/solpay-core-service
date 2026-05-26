package utils

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

// CreateJWTToken creates a signed JWT token string.
func CreateJWTToken(secretKey string, claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return tokenString, nil
}

// GenerateAccessToken creates a JWT access token for the given account.
func GenerateAccessToken(accountID uint) string {
	cfg := config.LoadConfig()

	hours, err := strconv.Atoi(cfg.JWT_EXPIRATION_HOURS)
	if err != nil {
		slog.Error("GenerateAccessToken: invalid JWT_EXPIRATION_HOURS", "error", err)
		return ""
	}

	claims := jwt.MapClaims{
		"account_id": accountID,
		"exp":        time.Now().Add(time.Hour * time.Duration(hours)).Unix(),
	}

	token, err := CreateJWTToken(cfg.SECRET_JWT, claims)
	if err != nil {
		slog.Error("GenerateAccessToken: failed to create JWT token", "error", err)
		return ""
	}
	return token
}
