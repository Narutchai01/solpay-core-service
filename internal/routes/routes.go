package routes

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RoutesConfig(app *fiber.App, db *gorm.DB) {
	var validate = validator.New()

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Example route group
	exampleGroup := v1.Group("/example")
	exampleRouteConfig := NewExampleRouteConfig(exampleGroup, db)
	exampleRouteConfig.Setup()

	// Account route group
	accountGroup := v1.Group("/accounts")
	accountRouteConfig := NewAccountRouteConfig(accountGroup, db, validate)
	accountRouteConfig.Setup()

	// Health route group
	healthGroup := v1.Group("/health")
	healthRouteConfig := NewHealthRouteConfig(healthGroup)
	healthRouteConfig.Setup()

	// Balance route group
	balanceGroup := v1.Group("/balances")
	balanceRouteConfig := NewBalanceRouteConfig(balanceGroup, db, validate)
	balanceRouteConfig.Setup()
}
