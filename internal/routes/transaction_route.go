package routes

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/handler"
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/infra/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TransactionRouteConfig struct {
	route    fiber.Router
	db       *gorm.DB
	validate *validator.Validate
}

func NewTransactionRouteConfig(route fiber.Router, db *gorm.DB, validate *validator.Validate) *TransactionRouteConfig {
	return &TransactionRouteConfig{
		route:    route,
		db:       db,
		validate: validate,
	}
}

func (trc *TransactionRouteConfig) Setup() {
	transactionRepository := repositories.NewGormTransactionRepository(trc.db)
	transactionService := services.NewTransactionService(transactionRepository)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	trc.route.Post("/", transactionHandler.CreateTransaction)
}
