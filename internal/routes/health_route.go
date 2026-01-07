package routes

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/handler"
	"github.com/gofiber/fiber/v2"
)

type HealthRouteConfig struct {
	route fiber.Router
}

func NewHealthRouteConfig(route fiber.Router) *HealthRouteConfig {
	return &HealthRouteConfig{
		route: route,
	}
}

func (hrc *HealthRouteConfig) Setup() {
	healthHandler := handler.NewHealthHandler()

	hrc.route.Get("/", healthHandler.HandleHealthCheck)
}
