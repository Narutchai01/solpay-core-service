package handler

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/dto/response"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/gofiber/fiber/v2"
)

// TransactionHandler defines HTTP handlers for transaction operations.
type TransactionHandler interface {
	CreateTransaction(c *fiber.Ctx) error
	GetTransactionByIDHandler(c *fiber.Ctx) error
	GetTransactionsHandler(c *fiber.Ctx) error
}

type transactionHandler struct {
	transactionService services.TransactionService
}

// NewTransactionHandler creates a new TransactionHandler.
func NewTransactionHandler(transactionService services.TransactionService) TransactionHandler {
	return &transactionHandler{
		transactionService: transactionService,
	}
}

func (h *transactionHandler) CreateTransaction(c *fiber.Ctx) error {
	var req request.CreateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		msg := utils.FormatValidationError(err)
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, msg, err))
	}

	transaction, err := h.transactionService.CreateTransaction(c.UserContext(), req)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	return utils.HandleResponse(c, response.FormatTransactionDTO(transaction), nil)
}

func (h *transactionHandler) GetTransactionByIDHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		msg := utils.FormatValidationError(err)
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, msg, err))
	}

	transaction, err := h.transactionService.GetTransactionByID(id)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	return utils.HandleResponse(c, response.FormatTransactionDTO(transaction), nil)
}

func (h *transactionHandler) GetTransactionsHandler(c *fiber.Ctx) error {
	var q request.TransactionQuery
	if err := c.QueryParser(&q); err != nil {
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, "invalid query parameters", err))
	}

	txs, total, err := h.transactionService.GetTransactions(0, q)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	formattedTxs := make([]*response.TransactionDTO, len(txs))
	for i, tx := range txs {
		formattedTxs[i] = response.FormatTransactionDTO(&tx)
	}

	return utils.HandleResponse(c, fiber.Map{
		"items":    formattedTxs,
		"total":    total,
		"page":     q.Page,
		"pageSize": q.PageSize,
	}, nil)
}
