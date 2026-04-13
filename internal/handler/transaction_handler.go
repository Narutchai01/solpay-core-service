package handler

import (
	"strconv"

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
	QueryTransactionSummaryHandler(c *fiber.Ctx) error
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
			"totalDeposit":  summary.TotalDeposit,
			"totalWithdraw": summary.TotalWithdraw,
		},
		"chartData": summary.ChartData,
	}, nil)
}
