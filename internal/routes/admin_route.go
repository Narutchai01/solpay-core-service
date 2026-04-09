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

type AdminRouteConfig struct {
	route fiber.Router
	db    *gorm.DB
}

func NewAdminRouteConfig(route fiber.Router, db *gorm.DB) *AdminRouteConfig {
	return &AdminRouteConfig{
		route: route,
		db:    db,
	}
}

func (arc *AdminRouteConfig) Setup() {
	cfg := config.LoadConfig()

	adminRepo := repositories.NewGormAdminRepository(arc.db)
	adminService := services.NewAdminService(adminRepo, cfg)
	adminHandler := handler.NewAdminHandler(adminService)

	arc.route.Post("/", adminHandler.CreateAdminHandler)
	arc.route.Post("/login", adminHandler.LoginAdminHandler)
	arc.route.Get("/me", middlewares.AuthRequired(), adminHandler.GetProfileHandler)
}
