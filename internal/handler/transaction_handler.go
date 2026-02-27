package handler

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/dto/response"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type TransactionHandler interface {
	CreateTransaction(c *fiber.Ctx) error
	GetTransactionByIDHandler(c *fiber.Ctx) error
}

type transactionHandler struct {
	transactionService services.TransactionService
}

func NewTransactionHandler(transactionService services.TransactionService) TransactionHandler {
	return &transactionHandler{
		transactionService: transactionService,
	}
}

func (h *transactionHandler) CreateTransaction(c *fiber.Ctx) error {
	var req request.CreateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	transaction, err := h.transactionService.CreateTransaction(c.UserContext(), req)
	if err != nil {
		return err
	}

	return utils.HandleResponse(c, transaction, nil)
}

func (h *transactionHandler) GetTransactionByIDHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		msg := utils.FormatValidationError(err)
		appErr := entities.NewAppError(entities.ErrTypeBadRequest, msg, err)
		return utils.HandleResponse(c, nil, appErr)
	}

	transaction, err := h.transactionService.GetTransactionByID(id)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	transactionDTO := response.FormaterTransactionDTO(transaction)

	return utils.HandleResponse(c, transactionDTO, nil)
}
