package handler

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type ExchangeRateHandler interface {
	GetExchangeRate(c *fiber.Ctx) error
}

type exchangeRateHandler struct {
	service services.ExchangeRateService
}

func NewExchangeRateHandler(service services.ExchangeRateService) ExchangeRateHandler {
	return &exchangeRateHandler{
		service: service,
	}
}

func (h *exchangeRateHandler) GetExchangeRate(c *fiber.Ctx) error {
	symbol := c.Query("sym")
	rates, err := h.service.GetExchangeRate(symbol)
	if err != nil {
		return utils.HandleResponse(c, nil, err)
	}

	return utils.HandleResponse(c, rates, nil, "Get exchange rates successfully")
}
