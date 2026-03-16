package middlewares

import (
	"strconv"
	"strings"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const invalidAccountIDClaimMsg = "Invalid account_id claim"

func unauthorized(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": msg,
	})
}

func parseBearerToken(authHeader string) (string, bool) {
	parts := strings.Fields(authHeader)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}
	return parts[1], true
}

func parseAccountIDClaim(claim any) (uint, bool) {
	switch value := claim.(type) {
	case float64:
		if value <= 0 {
			return 0, false
		}
		return uint(value), true
	case string:
		parsedValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil || parsedValue == 0 {
			return 0, false
		}
		return uint(parsedValue), true
	default:
		return 0, false
	}
}

func parseJWTClaims(tokenString, secret string) (jwt.MapClaims, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid signing method")
		}
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired())
	if err != nil || !token.Valid {
		return nil, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, false
	}

	return claims, true
}

// AuthRequired คือฟังก์ชัน Middleware สำหรับตรวจสอบ Token
func AuthRequired() fiber.Handler {

	cfg := config.LoadConfig()
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return unauthorized(c, "Authorization header is required")
		}

		tokenString, ok := parseBearerToken(authHeader)
		if !ok {
			return unauthorized(c, "Invalid token format")
		}

		claims, ok := parseJWTClaims(tokenString, cfg.SECRET_JWT)
		if !ok {
			return unauthorized(c, "Invalid or expired token")
		}

		accountClaim, ok := claims["account_id"]
		if !ok {
			return unauthorized(c, "account_id claim is required")
		}

		accountID, ok := parseAccountIDClaim(accountClaim)
		if !ok {
			return unauthorized(c, invalidAccountIDClaimMsg)
		}

		c.Locals("accountID", accountID)
		c.Locals("userID", int64(accountID))

		return c.Next()
	}
}
