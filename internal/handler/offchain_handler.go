package handler

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/dto/response"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type OffChainHandler interface {
	ConfirmOffChainHandler(c *fiber.Ctx) error
}

type offChainHandler struct {
	offchainService services.OffChainService
}

func NewOffChainHandler(offchainService services.OffChainService) OffChainHandler {
	return &offChainHandler{
		offchainService: offchainService,
	}
}

func (h *offChainHandler) ConfirmOffChainHandler(c *fiber.Ctx) error {
	var req request.OffChainRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, "invalid request body", err))
	}

	accountID, ok := utils.GetUserIDFromLocals(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	tx, err := h.offchainService.ComFirmOffchain(c.Context(), req, accountID)

	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	return utils.HandleResponse(c, response.FormatTransactionDTO(&tx), nil)

}
