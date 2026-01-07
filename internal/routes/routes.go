package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RoutesConfig(app *fiber.App, db *gorm.DB) {

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Example route group
	exampleGroup := v1.Group("/example")
	ExampleRoute(exampleGroup, db)

	// Account route group
	accountGroup := v1.Group("/accounts")
	AccountRoute(accountGroup, db)

	// Health route group
	healthGroup := v1.Group("/health")
	HealthRoute(healthGroup)
}
