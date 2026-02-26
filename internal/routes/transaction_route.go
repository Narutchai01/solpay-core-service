package routes

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/handler"
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/infra/rabbitmq"
	"github.com/Narutchai01/solpay-core-service/internal/infra/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type TransactionRouteConfig struct {
	route    fiber.Router
	db       *gorm.DB
	validate *validator.Validate
	channel  *amqp.Channel
}

func NewTransactionRouteConfig(route fiber.Router, db *gorm.DB, validate *validator.Validate, channel *amqp.Channel) *TransactionRouteConfig {
	return &TransactionRouteConfig{
		route:    route,
		db:       db,
		validate: validate,
		channel:  channel,
	}
}

func (trc *TransactionRouteConfig) Setup() {
	transactionRepository := repositories.NewGormTransactionRepository(trc.db)
	uowRepository := repositories.NewSqlUnitOfWork(trc.db)
	rabbitmqProducer := rabbitmq.NewProducer(trc.channel)
	transactionService := services.NewTransactionService(transactionRepository, uowRepository, rabbitmqProducer)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	trc.route.Post("/", transactionHandler.CreateTransaction)
	trc.route.Get("/:id", transactionHandler.GetTransactionByIDHandler)
}
