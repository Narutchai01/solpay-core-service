package routes

import (
	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/handler"
	"github.com/Narutchai01/solpay-core-service/internal/infra/rabbitmq"
	"github.com/Narutchai01/solpay-core-service/internal/infra/repositories"
	"github.com/Narutchai01/solpay-core-service/internal/middlewares"
	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type SwapRouteConfig struct {
	route   fiber.Router
	db      *gorm.DB
	channel *amqp.Channel
}

func NewSwapRouteConfig(route fiber.Router, db *gorm.DB, channel *amqp.Channel) *SwapRouteConfig {
	return &SwapRouteConfig{
		route:   route,
		db:      db,
		channel: channel,
	}
}

func (src *SwapRouteConfig) Setup() {
	cfg := config.LoadConfig()
	repo := repositories.NewSwapRepository(cfg.SWAP_SERVICE_URL)
	accountRepo := repositories.NewGormAccountRepository(src.db)
	transactionRepo := repositories.NewGormTransactionRepository(src.db)
	uow := repositories.NewSqlUnitOfWork(src.db)
	publisher := rabbitmq.NewPublisher(src.channel)

	service := services.NewSwapService(repo, accountRepo, transactionRepo, uow, publisher)
	swapHandler := handler.NewSwapHandler(service)

	src.route.Get("/quote", swapHandler.GetSwapQuote)
	src.route.Post("/swap", middlewares.AuthRequired(), swapHandler.BuildSwapUnsignedTransaction)
	src.route.Post("/execute", middlewares.AuthRequired(), swapHandler.ExecuteSwapTransaction)
}
