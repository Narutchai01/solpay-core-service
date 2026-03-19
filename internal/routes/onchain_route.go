package routes

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/handler"
	"github.com/Narutchai01/solpay-core-service/internal/infra/rabbitmq"
	"github.com/Narutchai01/solpay-core-service/internal/infra/repositories"
	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type OnchainRouteConfig struct {
	route   fiber.Router
	db      *gorm.DB
	channel *amqp.Channel
}

func NewOnchainRouteConfig(route fiber.Router, db *gorm.DB, channel *amqp.Channel) *OnchainRouteConfig {
	return &OnchainRouteConfig{
		route:   route,
		db:      db,
		channel: channel,
	}
}

func (trc *OnchainRouteConfig) Setup() {
	// สร้าง Transaction Repository, Quote Repository, Unit of Work, และ Service
	transactionRepo := repositories.NewGormTransactionRepository(trc.db)
	quoteRepo := repositories.NewGormQuoteRepository(trc.db)
	uow := repositories.NewSqlUnitOfWork(trc.db)
	pub := rabbitmq.NewPublisher(trc.channel)
	onchainService := services.NewOnchainService(transactionRepo, quoteRepo, uow, pub)

	// สร้าง Handler
	onchainHandler := handler.NewOnchainHandler(onchainService)

	// กำหนดเส้นทางสำหรับ On-Chain
	trc.route.Post("/confirm", onchainHandler.ConfirmOnchainHandler)
}
