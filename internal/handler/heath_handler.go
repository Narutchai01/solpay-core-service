package handler

import (
	"fmt"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/gofiber/fiber/v2"
)

// HealthHandler defines HTTP handlers for health checks.
type HealthHandler struct {
	cfg *config.Config
}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler(cfg *config.Config) *HealthHandler {
	return &HealthHandler{cfg: cfg}
}

func (h *HealthHandler) HandleHealthCheck(c *fiber.Ctx) error {
	msg := fmt.Sprintf("Service is running. Environment: %s, on port %s", h.cfg.Environment, h.cfg.APPPort)
	return utils.HandleResponse(c, nil, nil, msg)
}
