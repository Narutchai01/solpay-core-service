package handler

import (
	"log/slog"

	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/models/request"
	"github.com/Narutchai01/solpay-core-service/internal/models/response"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type AccountHandler interface {
	CreateAccountHandler(c *fiber.Ctx) error
}

type accountHandler struct {
	accountService services.AccountService
}

func NewAccountHandler(accountService services.AccountService) AccountHandler {
	return &accountHandler{
		accountService: accountService,
	}
}

func (h *accountHandler) CreateAccountHandler(c *fiber.Ctx) error {
	slog.Info("Received request to create account")
	var req request.CreateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewResponseModel(fiber.StatusBadRequest, "Invalid request body", nil, err.Error()))
	}

	account, err := h.accountService.CreateAccount(req)

	if err != nil {
		return utils.HandleError(c, err)
	}

	slog.Info("Account created successfully", "account", account)
	return c.Status(fiber.StatusCreated).JSON(response.NewResponseModel(fiber.StatusCreated, "Account created successfully", account, nil))
}
