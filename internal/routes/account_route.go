package routes

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/handler"
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/infra/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AccountRoute(route fiber.Router, db *gorm.DB, validate *validator.Validate) {
	accountRepo := repositories.NewGormAccountRepository(db)
	accountService := services.NewAccountService(accountRepo)
	accountHandler := handler.NewAccountHandler(accountService)

	route.Get("/", accountHandler.GetAccountsHandler)
	route.Post("/", accountHandler.CreateAccountHandler)
}
