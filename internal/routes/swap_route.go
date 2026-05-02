package routes

import (
	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/handler"
	"github.com/Narutchai01/solpay-core-service/internal/infra/repositories"
	"github.com/Narutchai01/solpay-core-service/internal/middlewares"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SwapRouteConfig struct {
	route fiber.Router
	db    *gorm.DB
}

func NewSwapRouteConfig(route fiber.Router, db *gorm.DB) *SwapRouteConfig {
	return &SwapRouteConfig{
		route: route,
		db:    db,
	}
}

func (src *SwapRouteConfig) Setup() {
	cfg := config.LoadConfig()
	repo := repositories.NewSwapRepository(cfg.SWAP_SERVICE_URL)
	accountRepo := repositories.NewGormAccountRepository(src.db)
	service := services.NewSwapService(repo, accountRepo)
	swapHandler := handler.NewSwapHandler(service)

	src.route.Get("/quote", swapHandler.GetSwapQuote)
	src.route.Post("/swap", middlewares.AuthRequired(), swapHandler.BuildSwapUnsignedTransaction)
}
