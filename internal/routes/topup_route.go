package routes

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/handler"
	"github.com/Narutchai01/solpay-core-service/internal/infra/repositories"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TopUpRouteConfig struct {
	route fiber.Router
	db    *gorm.DB
}

func NewTopUpRouteConfig(route fiber.Router, db *gorm.DB) *TopUpRouteConfig {
	return &TopUpRouteConfig{
		route: route,
		db:    db,
	}
}

func (trc *TopUpRouteConfig) Setup() {
	// สร้าง Transaction Repository, Quote Repository, Unit of Work, และ Service
	transactionRepo := repositories.NewGormTransactionRepository(trc.db)
	quoteRepo := repositories.NewGormQuoteRepository(trc.db)
	uow := repositories.NewSqlUnitOfWork(trc.db)
	transactionService := services.NewTransactionService(transactionRepo, uow, nil, nil)
	topUpService := services.NewTopUpService(transactionService, transactionRepo, quoteRepo, uow)

	// สร้าง Handler
	topUpHandler := handler.NewTopUpHandler(topUpService)

	// กำหนดเส้นทางสำหรับ Top-Up
	trc.route.Post("/confirm", topUpHandler.ConfirmTopUpHandler)
}
