package handler

import (
	"strconv"

	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type QuoteHandler interface {
	CreateQuoteHandler(c *fiber.Ctx) error
	GetQuoteByIDHandler(c *fiber.Ctx) error
}

type quoteHandler struct {
	quoteService services.QuoteService
}

func NewQuoteHandler(quoteService services.QuoteService) QuoteHandler {
	return &quoteHandler{
		quoteService: quoteService,
	}
}

func getUserIDFromLocals(c *fiber.Ctx) (int64, bool) {
	userID := c.Locals("userID")
	switch value := userID.(type) {
	case int64:
		return value, value > 0
	case uint:
		return int64(value), value > 0
	case int:
		return int64(value), value > 0
	case float64:
		return int64(value), value > 0
	case string:
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil || parsed <= 0 {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}

func (h *quoteHandler) CreateQuoteHandler(c *fiber.Ctx) error {
	var req request.CreateQuoteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	userID, ok := getUserIDFromLocals(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user context",
		})
	}

	quoteResp, err := h.quoteService.CreateQuote(req, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create quote",
		})
	}

	return c.JSON(quoteResp)
}

func (h *quoteHandler) GetQuoteByIDHandler(c *fiber.Ctx) error {
	quoteID := c.Params("id")
	quote, err := h.quoteService.GetQuoteByID(quoteID)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	return utils.HandleResponse(c, quote, nil, "receive quote by id successfully")
}
