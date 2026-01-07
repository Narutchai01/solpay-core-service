package handler

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type ExampleHandler struct {
	exampleService services.ExampleService
}

func NewExampleHandler(exampleService services.ExampleService) *ExampleHandler {
	return &ExampleHandler{
		exampleService: exampleService,
	}
}

func (h *ExampleHandler) HandleExampleGetById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		msg := utils.FormatValidationError(err)
		appErr := entities.NewAppError(entities.ErrTypeBadRequest, msg, err)
		return utils.HandleResponse(c, nil, appErr)
	}

	example, err := h.exampleService.GetExampleByID(id)

	return utils.HandleResponse(c, example, err)
}
