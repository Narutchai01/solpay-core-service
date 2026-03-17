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

type TopUpRouteConfig struct {
	route   fiber.Router
	db      *gorm.DB
	channel *amqp.Channel
}

func NewTopUpRouteConfig(route fiber.Router, db *gorm.DB, channel *amqp.Channel) *TopUpRouteConfig {
	return &TopUpRouteConfig{
		route:   route,
		db:      db,
		channel: channel,
	}
}

func (trc *TopUpRouteConfig) Setup() {
	// สร้าง Transaction Repository, Quote Repository, Unit of Work, และ Service
	transactionRepo := repositories.NewGormTransactionRepository(trc.db)
	quoteRepo := repositories.NewGormQuoteRepository(trc.db)
	uow := repositories.NewSqlUnitOfWork(trc.db)
	pub := rabbitmq.NewPublisher(trc.channel)
	topUpService := services.NewTopUpService(transactionRepo, quoteRepo, uow, pub)

	// สร้าง Handler
	topUpHandler := handler.NewTopUpHandler(topUpService)

	// กำหนดเส้นทางสำหรับ Top-Up
	trc.route.Post("/confirm", topUpHandler.ConfirmTopUpHandler)
}
