package routes

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/handler"
	"github.com/gofiber/fiber/v2"
)

func HealthRoute(route fiber.Router) {
	healthHandler := handler.NewHealthHandler()

	route.Get("/", healthHandler.HandleHealthCheck)
}
