package routes

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/handler"
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/infra/repositories"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ExampleRoute(route fiber.Router, db *gorm.DB) {

	exampleRepo := repositories.NewGormExampleRepository(db)
	exampleService := services.NewExampleService(exampleRepo)
	exampleHandler := handler.NewExampleHandler(exampleService)

	route.Get("/:id", exampleHandler.HandleExampleGetById)
}
