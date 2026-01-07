package handler

import (
	"fmt"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HandleHealthCheck(c *fiber.Ctx) error {

	cfg := config.LoadConfig()

	msg := fmt.Sprintf("Service is running. Environment: %s, on port %s", cfg.Environment, cfg.APPPort)

	return utils.HandleResponse(c, nil, nil, msg)
}
