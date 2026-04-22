package handler

import (
	"mime"
	"path/filepath"
	"strings"

	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	CreateUserHandler(c *fiber.Ctx) error
	ApproveUserHandler(c *fiber.Ctx) error
}

type userHandler struct {
	userService services.UserService
	validate    *validator.Validate
}

func NewUserHandler(userService services.UserService) UserHandler {
	return &userHandler{
		userService: userService,
		validate:    validator.New(),
	}
}

func (h *userHandler) CreateUserHandler(c *fiber.Ctx) error {
	frontCardImage, _ := c.FormFile("front_card_image")
	backCardImage, _ := c.FormFile("back_card_image")

	req := request.CreateUserRequest{
		IDCard:         c.FormValue("id_card"),
		FirstName:      c.FormValue("first_name"),
		LastName:       c.FormValue("last_name"),
		BirthDate:      c.FormValue("birth_date"),
		ExpireDate:     c.FormValue("expire_date"),
		FrontCardImage: frontCardImage,
		BackCardImage:  backCardImage,
	}

	if err := h.validate.Struct(&req); err != nil {
		msg := utils.FormatValidationError(err)
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, msg, err))
	}

	contentType := req.FrontCardImage.Header.Get("Content-Type")
	if contentType == "" {
		contentType = mime.TypeByExtension(filepath.Ext(req.FrontCardImage.Filename))
	}

	if !strings.HasPrefix(contentType, "image/") {
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, "front_card_image must be an image", entities.ErrBadRequest))
	}

	contentType = req.BackCardImage.Header.Get("Content-Type")
	if contentType == "" {
		contentType = mime.TypeByExtension(filepath.Ext(req.BackCardImage.Filename))
	}

	if !strings.HasPrefix(contentType, "image/") {
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, "back_card_image must be an image", entities.ErrBadRequest))
	}

	user, err := h.userService.CreateUser(&req)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	return utils.HandleResponse(c, user, nil)

}

func (h *userHandler) ApproveUserHandler(c *fiber.Ctx) error {
	var req request.ApprovalStatus
	if err := c.BodyParser(&req); err != nil {
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, "invalid request body", err))
	}

	if err := h.validate.Struct(&req); err != nil {
		msg := utils.FormatValidationError(err)
		return utils.HandleResponse(c, nil, entities.NewAppError(entities.ErrTypeBadRequest, msg, err))
	}

	user, err := h.userService.ApprovalStatus(&req)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	return utils.HandleResponse(c, user, nil)
}
