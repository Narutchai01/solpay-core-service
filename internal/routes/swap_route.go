package routes

import (
	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/handler"
	"github.com/Narutchai01/solpay-core-service/internal/infra/repositories"
	"github.com/gofiber/fiber/v2"
)

type SwapRouteConfig struct {
	route fiber.Router
}

func NewSwapRouteConfig(route fiber.Router) *SwapRouteConfig {
	return &SwapRouteConfig{
		route: route,
	}
}

func (src *SwapRouteConfig) Setup() {
	cfg := config.LoadConfig()
	repo := repositories.NewSwapRepository(cfg.SWAP_SERVICE_URL)
	service := services.NewSwapService(repo)
	swapHandler := handler.NewSwapHandler(service)

	src.route.Get("/quote", swapHandler.GetSwapQuote)
}
