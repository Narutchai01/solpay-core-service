package handler

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type CategoryHandler interface {
	GetCategoriesHandler(c *fiber.Ctx) error
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
