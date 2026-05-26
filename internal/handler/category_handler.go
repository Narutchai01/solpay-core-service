package handler

import (
	"strconv"

	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type CategoryHandler interface {
	GetCategoriesHandler(c *fiber.Ctx) error
	GetCategoryHandler(c *fiber.Ctx) error
}

type categoryHandler struct {
	categoryService services.CategoryService
}

func NewCategoryHandler(categoryService services.CategoryService) CategoryHandler {
	return &categoryHandler{categoryService: categoryService}
}

func (h *categoryHandler) GetCategoriesHandler(c *fiber.Ctx) error {
	categories, err := h.categoryService.GetCategories()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get categories")
	}
	return utils.HandleResponse(c, categories, nil)
}

func (h *categoryHandler) GetCategoryHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid category id")
	}

	category, err := h.categoryService.GetCategory(id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "category not found")
	}

	return utils.HandleResponse(c, category, nil)
}
