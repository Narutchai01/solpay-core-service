package routes

import (
	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/handler"
	"github.com/Narutchai01/solpay-core-service/internal/infra/bithub"
	"github.com/gofiber/fiber/v2"
)

type ExchangeRateRouteConfig struct {
	route fiber.Router
	cfg   *config.Config
}

func NewExchangeRateRouteConfig(route fiber.Router, cfg *config.Config) *ExchangeRateRouteConfig {
	return &ExchangeRateRouteConfig{
		route: route,
		cfg:   cfg,
	}
}

func (erc *ExchangeRateRouteConfig) Setup() {
	repo := bithub.NewExchangeRateRepository(erc.cfg.EXCHANGE_RATE_URL)
	service := services.NewExchangeRateService(repo)
	exchangeRateHandler := handler.NewExchangeRateHandler(service)

	erc.route.Get("/", exchangeRateHandler.GetExchangeRate)
}
