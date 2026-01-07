package utils

import (
	"github.com/Narutchai01/solpay-core-service/internal/dto/response"
	"github.com/gofiber/fiber/v2"
)

func HandleSuccess(c *fiber.Ctx, status int, msg string, data interface{}) error {
	return c.Status(status).JSON(response.FormaterResponseDTO(status, msg, data, nil))
}
