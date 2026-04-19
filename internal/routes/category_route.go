package routes

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/handler"
	"github.com/Narutchai01/solpay-core-service/internal/infra/repositories"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CategoryRouteConfig struct {
	route fiber.Router
	db    *gorm.DB
}

func NewCategoryRouteConfig(route fiber.Router, db *gorm.DB) *CategoryRouteConfig {
	return &CategoryRouteConfig{
		route: route,
		db:    db,
	}
}

func (crc *CategoryRouteConfig) Setup() {
	categoryRepo := repositories.NewCategoryRepository(crc.db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	crc.route.Get("/", categoryHandler.GetCategoriesHandler)
}
