package handler

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type TopUpHandler interface {
	ConfirmTopUpHandler(c *fiber.Ctx) error
}

type topUpHandler struct {
	topUpService services.TopUpService
}

func NewTopUpHandler(topUpService services.TopUpService) TopUpHandler {
	return &topUpHandler{
		topUpService: topUpService,
	}
}

func (h *topUpHandler) ConfirmTopUpHandler(c *fiber.Ctx) error {
	var req request.TopUpRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	accountID, ok := utils.GetUserIDFromLocals(c)
	if !ok {
		return utils.HandleResponse(c, nil, fiber.NewError(fiber.StatusUnauthorized, "unauthorized"))
	}

	tx, err := h.topUpService.ComfirmTopUp(c.Context(), req, accountID)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	return utils.HandleResponse(c, tx, nil)
}
