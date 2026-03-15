package handler

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/gofiber/fiber/v2"
)

type QuoteHandler interface {
	CreateQuoteHandler(c *fiber.Ctx) error
}

type quoteHandler struct {
	quoteService services.QuoteService
}

func NewQuoteHandler(quoteService services.QuoteService) QuoteHandler {
	return &quoteHandler{
		quoteService: quoteService,
	}
}

func (h *quoteHandler) CreateQuoteHandler(c *fiber.Ctx) error {
	var req request.CreateQuoteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	quoteResp, err := h.quoteService.CreateQuote(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create quote",
		})
	}

	return c.JSON(quoteResp)
}
