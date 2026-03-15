package routes

import (
	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

// RoutesConfig sets up all API route groups.
func RoutesConfig(app *fiber.App, db *gorm.DB, channel *amqp.Channel, cfg *config.Config) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	NewExampleRouteConfig(v1.Group("/example"), db).Setup()
	NewAccountRouteConfig(v1.Group("/accounts"), db).Setup()
	NewHealthRouteConfig(v1.Group("/health"), cfg).Setup()
	NewBalanceRouteConfig(v1.Group("/balances"), db).Setup()
	NewTransactionRouteConfig(v1.Group("/transactions"), db, channel).Setup()
	NewQuoteRouteConfig(v1.Group("/quotes"), db).Setup()
	NewTopUpRouteConfig(v1.Group("/topup"), db).Setup()
}
