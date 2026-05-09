package handler

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type SwapHandler interface {
	GetSwapQuote(c *fiber.Ctx) error
	BuildSwapUnsignedTransaction(c *fiber.Ctx) error
	BuildSwapInstruction(c *fiber.Ctx) error
	ExecuteSwapTransaction(c *fiber.Ctx) error
}

type swapHandler struct {
	service services.SwapService
}

func NewSwapHandler(service services.SwapService) SwapHandler {
	return &swapHandler{service: service}
}

func (h *swapHandler) GetSwapQuote(c *fiber.Ctx) error {
	var query request.SwapQuoteRequest
	if err := c.QueryParser(&query); err != nil {
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, "invalid query parameters", err))
	}

	resp, err := h.service.GetSwapQuote(query)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	return utils.HandleResponse(c, resp, nil, "get swap quote successfully")
}

func (h *swapHandler) BuildSwapUnsignedTransaction(c *fiber.Ctx) error {
	userID, ok := utils.GetUserIDFromLocals(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user context",
		})
	}

	var req request.SwapUnsignedTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, "invalid request body", err))
	}

	resp, err := h.service.BuildSwapUnsignedTransaction(req, userID)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	return utils.HandleResponse(c, resp, nil, "build swap unsigned transaction successfully")
}

func (h *swapHandler) BuildSwapInstruction(c *fiber.Ctx) error {
	userID, ok := utils.GetUserIDFromLocals(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user context",
		})
	}

	var req request.SwapUnsignedTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, "invalid request body", err))
	}

	resp, err := h.service.BuildSwapInstruction(req, userID)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	return utils.HandleResponse(c, resp, nil, "build swap instruction successfully")
}

func (h *swapHandler) ExecuteSwapTransaction(c *fiber.Ctx) error {
	userID, ok := utils.GetUserIDFromLocals(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user context",
		})
	}

	var req request.ExecuteSwapTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, "invalid request body", err))
	}

	resp, err := h.service.ExecuteSwapTransaction(c.Context(), req, userID)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	return utils.HandleResponse(c, resp, nil, "execute swap transaction successfully")
}
