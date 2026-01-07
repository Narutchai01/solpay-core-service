package routes

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/handler"
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/infra/repositories"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ExampleRouteConfig struct {
	route fiber.Router
	db    *gorm.DB
}

func NewExampleRouteConfig(route fiber.Router, db *gorm.DB) *ExampleRouteConfig {
	return &ExampleRouteConfig{
		route: route,
		db:    db,
	}
}

func (erc *ExampleRouteConfig) Setup() {
	exampleRepo := repositories.NewGormExampleRepository(erc.db)
	exampleService := services.NewExampleService(exampleRepo)
	exampleHandler := handler.NewExampleHandler(exampleService)

	erc.route.Get("/:id", exampleHandler.HandleExampleGetById)
}
