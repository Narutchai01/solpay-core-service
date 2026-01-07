package handler

import (
	"fmt"
	"log/slog"

	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/dto/response"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AccountHandler interface {
	CreateAccountHandler(c *fiber.Ctx) error
	GetAccountsHandler(c *fiber.Ctx) error
}

type accountHandler struct {
	accountService services.AccountService
	validate       *validator.Validate
}

func NewAccountHandler(accountService services.AccountService) AccountHandler {
	return &accountHandler{
		accountService: accountService,
		validate:       validator.New(),
	}
}

// NOTE: CreateAccountHandler handles the creation of a new account
func (h *accountHandler) CreateAccountHandler(c *fiber.Ctx) error {
	// NOTE: Parse and validate the request body
	var req request.CreateAccountRequest
	// NOTE: Handle parsing and validation errors
	if err := c.BodyParser(&req); err != nil {
		msg := utils.FormatValidationError(err)
		// return utils.HandleError(c, entities.NewAppError(entities.ErrTypeBadRequest, msg, err))
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, msg, err))
	}

	// NOTE: Validate the request struct
	if err := h.validate.Struct(&req); err != nil {
		msg := utils.FormatValidationError(err)
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, msg, err))
	}

	// NOTE: Call the service to create the account
	account, err := h.accountService.CreateAccount(req)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	// NOTE: define success message
	msg := fmt.Sprintf("Account %d created successfully", account.ID)

	slog.Info(msg)
	return utils.HandleResponse(c, account, nil, msg)
}

func (h *accountHandler) GetAccountsHandler(c *fiber.Ctx) error {

	// NOTE: Get pagination query parameters
	// FIXME: fix logic to use request struct
	var req request.GetAccountsRequest
	if err := c.QueryParser(&req); err != nil {
		msg := utils.FormatValidationError(err)
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, msg, err))
	}

	if err := h.validate.Struct(&req); err != nil {
		msg := utils.FormatValidationError(err)
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, msg, err))
	}

	page, limit := req.Page, req.Limit

	// NOTE: Call the service to get accounts
	accounts, err := h.accountService.GetAccounts(page, limit)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	pagination := response.FormaterPaginationResponseDTO(100, page, accounts) // FIXME: fix total and totalPages

	// NOTE: define success message
	msg := fmt.Sprintf("Retrieved %d accounts successfully", len(accounts))

	slog.Info(msg)
	return utils.HandleResponse(c, pagination, nil, msg)
}
