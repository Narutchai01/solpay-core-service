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

// TransactionRouteConfig configures transaction-related routes.
type TransactionRouteConfig struct {
	route   fiber.Router
	db      *gorm.DB
	channel *amqp.Channel
}

// NewTransactionRouteConfig creates a new TransactionRouteConfig.
func NewTransactionRouteConfig(route fiber.Router, db *gorm.DB, channel *amqp.Channel) *TransactionRouteConfig {
	return &TransactionRouteConfig{
		route:   route,
		db:      db,
		channel: channel,
	}
}

func (trc *TransactionRouteConfig) Setup() {
	transactionRepo := repositories.NewGormTransactionRepository(trc.db)
	uow := repositories.NewSqlUnitOfWork(trc.db)
	pub := rabbitmq.NewPublisher(trc.channel)
	transactionService := services.NewTransactionService(transactionRepo, uow, pub, nil)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	trc.route.Post("/", transactionHandler.CreateTransaction)
	trc.route.Get("/", transactionHandler.GetTransactionsHandler)
	trc.route.Get("/summary", transactionHandler.QueryTransactionSummaryHandler)
	trc.route.Get("/me", middlewares.AuthRequired(), transactionHandler.GetTransactionsByAccountIDHandler)
	trc.route.Get("/:uuid", transactionHandler.GetTransactionByUUIDHandler)
}
