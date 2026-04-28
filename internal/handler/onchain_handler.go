package handler

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/dto/response"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type OnchainHandler interface {
	ConfirmOnchainHandler(c *fiber.Ctx) error
}

type onchainHandler struct {
	onchainService services.OnchainService
}

func NewOnchainHandler(onchainService services.OnchainService) OnchainHandler {
	return &onchainHandler{
		onchainService: onchainService,
	}
}

func (h *onchainHandler) ConfirmOnchainHandler(c *fiber.Ctx) error {
	var req request.TopUpRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.CategoryID == 0 {
		req.CategoryID = 1
	}

	accountID, ok := utils.GetUserIDFromLocals(c)
	if !ok {
		return utils.HandleResponse(c, nil, fiber.NewError(fiber.StatusUnauthorized, "unauthorized"))
	}

	tx, err := h.onchainService.ComfirmOnchain(c.Context(), req, accountID)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	return utils.HandleResponse(c, response.FormatTransactionDTO(&tx), nil)
}
