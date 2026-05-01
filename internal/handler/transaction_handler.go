package handler

import (
	"strconv"
	"strings"

	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/dto/response"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// TransactionHandler defines HTTP handlers for transaction operations.
type TransactionHandler interface {
	CreateTransaction(c *fiber.Ctx) error
	GetTransactionByUUIDHandler(c *fiber.Ctx) error
	GetTransactionsHandler(c *fiber.Ctx) error
	QueryTransactionSummaryHandler(c *fiber.Ctx) error
	GetTransactionsByAccountIDHandler(c *fiber.Ctx) error
	GetSpendingSummaryHandler(c *fiber.Ctx) error
	GetMonthlySpendingSummaryHandler(c *fiber.Ctx) error
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

func (h *transactionHandler) GetTransactionByUUIDHandler(c *fiber.Ctx) error {
	txUUIDStr := strings.TrimSpace(c.Params("uuid"))

	txUUID, err := uuid.Parse(txUUIDStr)
	if err != nil {
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, "Invalid input data", err))
	}

	transaction, err := h.transactionService.GetTransactionByUUID(txUUID)
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

	txs, total, err := h.transactionService.GetTransactions(q, nil)
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

func (h *transactionHandler) QueryTransactionSummaryHandler(c *fiber.Ctx) error {
	query := new(request.QueryTransactionSummaryRequest)

	if err := c.QueryParser(query); err != nil {
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, "invalid query parameters", err))
	}

	month, err := strconv.Atoi(query.Month)
	if err != nil || month < 1 || month > 12 {
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, "month must be between 1 and 12", nil))
	}

	year, err := strconv.Atoi(query.Year)
	if err != nil || year < 2000 {
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, "invalid year", nil))
	}

	summary, err := h.transactionService.QueryTransactionSummary(c.UserContext(), month, year)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	return utils.HandleResponse(c, fiber.Map{
		"summary": fiber.Map{
			"totalDeposit":        summary.TotalDeposit,
			"totalWithdraw":       summary.TotalWithdraw,
			"totalFee":            summary.TotalFee,
			"totalCompletedCount": summary.TotalCompletedCount,
		},
		"chartData": summary.ChartData,
	}, nil)
}

func (h *transactionHandler) GetTransactionsByAccountIDHandler(c *fiber.Ctx) error {

	accountID, ok := utils.GetUserIDFromLocals(c)
	if !ok {
		return utils.HandleResponse(c, nil, fiber.NewError(fiber.StatusUnauthorized, "unauthorized"))
	}

	var q request.TransactionQuery
	if err := c.QueryParser(&q); err != nil {
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, "invalid query parameters", err))
	}

	txs, total, err := h.transactionService.GetTransactions(q, &accountID)
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

func (h *transactionHandler) GetSpendingSummaryHandler(c *fiber.Ctx) error {
	accountID, ok := utils.GetUserIDFromLocals(c)
	if !ok {
		return utils.HandleResponse(c, nil, fiber.NewError(fiber.StatusUnauthorized, "unauthorized"))
	}

	summary, err := h.transactionService.GetSpendingSummary(c.UserContext(), accountID)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	return utils.HandleResponse(c, summary, nil)
}

func (h *transactionHandler) GetMonthlySpendingSummaryHandler(c *fiber.Ctx) error {
	accountID, ok := utils.GetUserIDFromLocals(c)
	if !ok {
		return utils.HandleResponse(c, nil, fiber.NewError(fiber.StatusUnauthorized, "unauthorized"))
	}

	summary, err := h.transactionService.GetMonthlySpendingSummary(c.UserContext(), accountID)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	return utils.HandleResponse(c, summary, nil)
}
