package routes

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/handler"
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/infra/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BalanceRouteConfig struct {
	route    fiber.Router
	db       *gorm.DB
	validate *validator.Validate
}

func NewBalanceRouteConfig(route fiber.Router, db *gorm.DB, validate *validator.Validate) *BalanceRouteConfig {
	return &BalanceRouteConfig{
		route:    route,
		db:       db,
		validate: validate,
	}
}

func (brc *BalanceRouteConfig) Setup() {
	balanceRepository := repositories.NewGormBalanceRepository(brc.db)
	uowRepository := repositories.NewSqlUnitOfWork(brc.db)
	balanceService := services.NewBalanceService(balanceRepository, uowRepository)
	balanceHandler := handler.NewBalanceHandler(balanceService)

	brc.route.Get("/", balanceHandler.GetBalancesHandler)
	brc.route.Get("/:id", balanceHandler.GetBalanceByIDHandler)
}
