package utils

import (
	"errors"
	"log/slog"

	"github.com/Narutchai01/solpay-core-service/internal/dto/response"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/gofiber/fiber/v2"
)

func handleError(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "Internal Server Error"

	var appErr *entities.AppError
	if errors.As(err, &appErr) {
		switch appErr.Type {
		case entities.ErrTypeConflict:
			code = fiber.StatusConflict
			msg = appErr.Message
			slog.Warn("Business Error", "msg", msg, "details", appErr.Err)
		case entities.ErrTypeNotFound:
			code = fiber.StatusNotFound
			msg = appErr.Message
			slog.Warn("Not Found", "msg", msg)
		case entities.ErrTypeBadRequest:
			code = fiber.StatusBadRequest
			msg = appErr.Message
			slog.Warn("Bad Request", "msg", msg)
		default:
			slog.Error("System Error", "msg", msg, "error", appErr.Err)
		}
	} else {
		slog.Error("Unknown Error", "error", err)
	}

	return c.Status(code).JSON(response.FormaterResponseDTO(code, msg, nil, msg))
}

func handleSuccess(c *fiber.Ctx, status int, msg string, data interface{}) error {
	return c.Status(status).JSON(response.FormaterResponseDTO(status, msg, data, nil))
}

func HandleResponse(c *fiber.Ctx, data interface{}, err error, message ...string) error {
	if err != nil {
		return handleError(c, err)
	}
	code := fiber.StatusOK
	if c.Method() == fiber.MethodPost {
		code = fiber.StatusCreated
	}

	msg := "Success"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	return handleSuccess(c, code, msg, data)

}
