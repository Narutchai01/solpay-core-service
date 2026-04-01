package routes

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/handler"
	"github.com/Narutchai01/solpay-core-service/internal/infra/rabbitmq"
	"github.com/Narutchai01/solpay-core-service/internal/infra/repositories"
	"github.com/Narutchai01/solpay-core-service/internal/middlewares"
	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type OffChainRouteConfig struct {
	route   fiber.Router
	db      *gorm.DB
	channel *amqp.Channel
}

func NewOffChainRouteConfig(route fiber.Router, db *gorm.DB, channel *amqp.Channel) *OffChainRouteConfig {
	return &OffChainRouteConfig{
		route:   route,
		db:      db,
		channel: channel,
	}
}

func (orc *OffChainRouteConfig) Setup() {
	// สร้าง Transaction Repository, Quote Repository, Unit of Work, และ Service
	transactionRepo := repositories.NewGormTransactionRepository(orc.db)
	quoteRepo := repositories.NewGormQuoteRepository(orc.db)
	uow := repositories.NewSqlUnitOfWork(orc.db)
	pub := rabbitmq.NewPublisher(orc.channel)
	offchainService := services.NewOffChainService(transactionRepo, uow, pub, quoteRepo)

	// สร้าง Handler
	offchainHandler := handler.NewOffChainHandler(offchainService)

	// กำหนดเส้นทางสำหรับ Off-Chain
	orc.route.Post("/confirm", middlewares.AuthRequired(), offchainHandler.ConfirmOffChainHandler)
}
