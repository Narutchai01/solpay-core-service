package handler

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type QuoteHandler interface {
	CreateQuoteHandler(c *fiber.Ctx) error
	GetQuoteByIDHandler(c *fiber.Ctx) error
	ConFirmQuoteHandler(c *fiber.Ctx) error
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
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, "invalid request body", err))
	}

	userID, ok := utils.GetUserIDFromLocals(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user context",
		})
	}

	quoteResp, err := h.quoteService.CreateQuote(req, userID)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	return utils.HandleResponse(c, quoteResp, nil, "create quote successfully")
}

func (h *quoteHandler) GetQuoteByIDHandler(c *fiber.Ctx) error {
	quoteID := c.Params("id")
	quote, err := h.quoteService.GetQuoteByID(quoteID)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	return utils.HandleResponse(c, quote, nil, "receive quote by id successfully")
}

func (h *quoteHandler) ConFirmQuoteHandler(c *fiber.Ctx) error {
	quoteID := c.Params("id")

	userID, ok := utils.GetUserIDFromLocals(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user context",
		})
	}

	txHash, err := h.quoteService.ConFirmQuote(quoteID, userID)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	return utils.HandleResponse(c, fiber.Map{"tx_hash": txHash}, nil, "confirm quote successfully")
}
