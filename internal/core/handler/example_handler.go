package handler

import (
	service "github.com/Narutchai01/solpay-core-service/internal/core/ports/services"
	"github.com/gofiber/fiber/v2"
)

type ExampleHandler struct {
	exampleService service.ExampleService
}

func NewExampleHandler(exampleService service.ExampleService) *ExampleHandler {
	return &ExampleHandler{
		exampleService: exampleService,
	}
}

func (h *ExampleHandler) HandleExampleGetById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID parameter",
		})
	}
	example, err := h.exampleService.GetExampleByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve example",
		})
	}
	return c.Status(fiber.StatusOK).JSON(example)
}
