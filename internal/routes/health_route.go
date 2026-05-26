package routes

import (
	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/handler"
	"github.com/gofiber/fiber/v2"
)

// HealthRouteConfig configures health-check routes.
type HealthRouteConfig struct {
	route fiber.Router
	cfg   *config.Config
}

// NewHealthRouteConfig creates a new HealthRouteConfig.
func NewHealthRouteConfig(route fiber.Router, cfg *config.Config) *HealthRouteConfig {
	return &HealthRouteConfig{
		route: route,
		cfg:   cfg,
	}
}

func (hrc *HealthRouteConfig) Setup() {
	healthHandler := handler.NewHealthHandler(hrc.cfg)
	hrc.route.Get("/", healthHandler.HandleHealthCheck)
}
