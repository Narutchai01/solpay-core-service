package routes

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/handler"
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/infra/repositories"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AccountRoute(route fiber.Router, db *gorm.DB) {
	accountRepo := repositories.NewGormAccountRepository(db)
	accountService := services.NewAccountService(accountRepo)
	accountHandler := handler.NewAccountHandler(accountService)

	route.Post("/", accountHandler.CreateAccountHandler)
}
