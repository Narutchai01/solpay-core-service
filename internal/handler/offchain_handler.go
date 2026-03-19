package handler

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
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
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	tx, err := h.offchainService.ComFirmOffchain(c.Context(), req)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to confirm off-chain transaction")
	}

	return utils.HandleResponse(c, tx, nil)

}
