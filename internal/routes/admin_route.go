package routes

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/handler"
	"github.com/Narutchai01/solpay-core-service/internal/infra/repositories"
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
	adminRepo := repositories.NewGormAdminRepository(arc.db)
	uow := repositories.NewSqlUnitOfWork(arc.db)
	adminService := services.NewAdminService(adminRepo, uow)
	adminHandler := handler.NewAdminHandler(adminService)

	arc.route.Post("/", adminHandler.CreateAdminHandler)
}
