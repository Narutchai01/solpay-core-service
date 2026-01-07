package utils

import (
	"github.com/Narutchai01/solpay-core-service/internal/models/response"
	"github.com/gofiber/fiber/v2"
)

func HandleSuccess(c *fiber.Ctx, status int, msg string, data interface{}) error {
	return c.Status(status).JSON(response.NewResponseModel(status, msg, data, nil))
}
