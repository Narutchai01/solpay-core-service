package handler

import (
	"fmt"

	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/dto/response"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// BalanceHandler defines HTTP handlers for balance operations.
type BalanceHandler interface {
	GetBalancesHandler(c *fiber.Ctx) error
	GetBalanceByIDHandler(c *fiber.Ctx) error
}

type balanceHandler struct {
	balanceService services.BalanceService
	validate       *validator.Validate
}

// NewBalanceHandler creates a new BalanceHandler.
func NewBalanceHandler(balanceService services.BalanceService) BalanceHandler {
	return &balanceHandler{
		balanceService: balanceService,
		validate:       validator.New(),
	}
}

func (h *balanceHandler) GetBalancesHandler(c *fiber.Ctx) error {
	var req request.GetBalancesRequest
	if err := c.QueryParser(&req); err != nil {
		msg := utils.FormatValidationError(err)
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, msg, err))
	}

	if err := h.validate.Struct(&req); err != nil {
		msg := utils.FormatValidationError(err)
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, msg, err))
	}

	balances, total, err := h.balanceService.GetBalances(req.Page, req.Limit)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	pagination := response.FormatPaginationResponseDTO(int(total), req.Page, response.FormatBalanceDTOs(balances))
	msg := fmt.Sprintf("Retrieved %d balances successfully", len(balances))

	return utils.HandleResponse(c, pagination, nil, msg)
}

func (h *balanceHandler) GetBalanceByIDHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		msg := utils.FormatValidationError(err)
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, msg, err))
	}

	balance, err := h.balanceService.GetBalanceByID(id)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	return utils.HandleResponse(c, response.FormatBalanceDTO(balance), nil)
}
