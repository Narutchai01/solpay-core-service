package handler

import (
	"fmt"
	"log/slog"

	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/models/request"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AccountHandler interface {
	CreateAccountHandler(c *fiber.Ctx) error
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

func (h *accountHandler) CreateAccountHandler(c *fiber.Ctx) error {
	var req request.CreateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		msg := utils.FormatValidationError(err)
		return utils.HandleError(c, entities.NewAppError(entities.ErrTypeBadRequest, msg, err))
	}

	if err := h.validate.Struct(&req); err != nil {
		msg := utils.FormatValidationError(err)
		return utils.HandleError(c, entities.NewAppError(entities.ErrTypeBadRequest, msg, err))
	}

	account, err := h.accountService.CreateAccount(req)
	if err != nil {
		return utils.HandleError(c, err)
	}

	msg := fmt.Sprintf("Account %d created successfully", account.ID)

	slog.Info(msg)
	return utils.HandleSuccess(c, fiber.StatusCreated, msg, account)
}
