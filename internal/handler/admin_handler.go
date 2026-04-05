package handler

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/dto/response"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AdminHandler interface {
	CreateAdminHandler(c *fiber.Ctx) error
}

type adminHandler struct {
	adminService services.AdminService
	validate     *validator.Validate
}

func NewAdminHandler(adminService services.AdminService) AdminHandler {
	return &adminHandler{
		adminService: adminService,
		validate:     validator.New(),
	}
}

func (h *adminHandler) CreateAdminHandler(c *fiber.Ctx) error {
	var req request.CreateAdminRequest
	if err := c.BodyParser(&req); err != nil {
		msg := utils.FormatValidationError(err)
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, msg, err))
	}

	if err := h.validate.Struct(&req); err != nil {
		msg := utils.FormatValidationError(err)
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, msg, err))
	}

	admin, err := h.adminService.CreateAdmin(c.UserContext(), req)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}
	resDTO := response.FormatAdminDTO(admin)
	return utils.HandleResponse(c, resDTO, nil)
}
