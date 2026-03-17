package handler

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
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

	return nil
}
