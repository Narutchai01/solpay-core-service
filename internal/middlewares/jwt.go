package middlewares

import (
	"strings"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return fiber.ErrUnauthorized
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return fiber.ErrUnauthorized
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.SECRET_JWT), nil
		})

		if err != nil || !token.Valid {
			return fiber.ErrUnauthorized
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return fiber.ErrUnauthorized
		}

		sub, ok := claims["sub"].(string)
		if !ok {
			return fiber.ErrUnauthorized
		}

		c.Locals("user_id", sub)

		return c.Next()
	}
}
