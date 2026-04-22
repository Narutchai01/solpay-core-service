package routes

import (
	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/handler"
	"github.com/Narutchai01/solpay-core-service/internal/infra/repositories"
	"github.com/Narutchai01/solpay-core-service/internal/infra/supabase"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserRouteConfig struct {
	route fiber.Router
	db    *gorm.DB
	cfg   *config.Config
}

func NewUserRouteConfig(route fiber.Router, db *gorm.DB, cfg *config.Config) *UserRouteConfig {
	return &UserRouteConfig{
		route: route,
		db:    db,
		cfg:   cfg,
	}
}

func (urc *UserRouteConfig) Setup() {
	userRepo := repositories.NewGormUserRepository(urc.db)
	storage := supabase.NewSupabaseStorage(urc.cfg.SUPABASE_PRIVATE_KEY, urc.cfg.SUPABASE_URL)
	userService := services.NewUserService(userRepo, storage)
	userHandler := handler.NewUserHandler(userService)

	urc.route.Post("/", userHandler.CreateUserHandler)
}
