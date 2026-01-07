package utils

import (
	"errors"
	"log/slog"

	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/models/response"
	"github.com/gofiber/fiber/v2"
)

func HandleError(c *fiber.Ctx, err error) error {
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

	return c.Status(code).JSON(response.NewResponseModel(code, msg, nil, msg))
}
