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

// AccountHandler defines HTTP handlers for account operations.
type AccountHandler interface {
	CreateAccountHandler(c *fiber.Ctx) error
	GetAccountsHandler(c *fiber.Ctx) error
	GetAccountByIDHandler(c *fiber.Ctx) error
	GetAccountByProfileHandler(c *fiber.Ctx) error
}

type accountHandler struct {
	accountService services.AccountService
	validate       *validator.Validate
}

// NewAccountHandler creates a new AccountHandler.
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
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, msg, err))
	}

	if err := h.validate.Struct(&req); err != nil {
		msg := utils.FormatValidationError(err)
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, msg, err))
	}

	account, err := h.accountService.CreateAccount(c.UserContext(), req)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	msg := fmt.Sprintf("Account %d", account.ID)

	accessToken := utils.GenerateAccessToken(account.ID)
	tokenDTO := response.TokenDTO{
		AccessToken:  accessToken,
		RefreshToken: "", // TODO: implement refresh token
	}

	return utils.HandleResponse(c, tokenDTO, nil, msg)
}

func (h *accountHandler) GetAccountsHandler(c *fiber.Ctx) error {
	var req request.GetAccountsRequest
	if err := c.QueryParser(&req); err != nil {
		msg := utils.FormatValidationError(err)
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, msg, err))
	}

	if err := h.validate.Struct(&req); err != nil {
		msg := utils.FormatValidationError(err)
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, msg, err))
	}

	accounts, total, err := h.accountService.GetAccounts(req.Page, req.Limit)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	pagination := response.FormatPaginationResponseDTO(int(total), req.Page, response.FormatAccountDTOs(accounts))
	msg := fmt.Sprintf("Retrieved %d accounts successfully", len(accounts))

	return utils.HandleResponse(c, pagination, nil, msg)
}

func (h *accountHandler) GetAccountByIDHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		msg := utils.FormatValidationError(err)
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, msg, err))
	}

	account, err := h.accountService.GetAccountByID(id)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	return utils.HandleResponse(c, response.FormatAccountDTO(account), nil)
}

func (h *accountHandler) GetAccountByProfileHandler(c *fiber.Ctx) error {
	accountID, ok := utils.GetUserIDFromLocals(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user context",
		})
	}

	account, err := h.accountService.GetAccountByID(int(accountID))
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	return utils.HandleResponse(c, response.FormatAccountDTO(account), nil)
}
