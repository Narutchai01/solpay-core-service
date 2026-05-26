package handler

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/gofiber/fiber/v2"
)

// ExampleHandler defines HTTP handlers for example operations.
type ExampleHandler struct {
	exampleService services.ExampleService
}

// NewExampleHandler creates a new ExampleHandler.
func NewExampleHandler(exampleService services.ExampleService) *ExampleHandler {
	return &ExampleHandler{
		exampleService: exampleService,
	}
}

func (h *ExampleHandler) HandleExampleGetById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		msg := utils.FormatValidationError(err)
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, msg, err))
	}

	example, err := h.exampleService.GetExampleByID(id)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	return utils.HandleResponse(c, example, nil)
}
