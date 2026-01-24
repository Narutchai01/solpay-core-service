package handler

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type TransactionHandler interface {
	CreateTransaction(c *fiber.Ctx) error
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
