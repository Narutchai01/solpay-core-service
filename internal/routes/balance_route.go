package routes

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/handler"
	"github.com/Narutchai01/solpay-core-service/internal/infra/repositories"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// BalanceRouteConfig configures balance-related routes.
type BalanceRouteConfig struct {
	route fiber.Router
	db    *gorm.DB
}

// NewBalanceRouteConfig creates a new BalanceRouteConfig.
func NewBalanceRouteConfig(route fiber.Router, db *gorm.DB) *BalanceRouteConfig {
	return &BalanceRouteConfig{
		route: route,
		db:    db,
	}
}

func (brc *BalanceRouteConfig) Setup() {
	balanceRepo := repositories.NewGormBalanceRepository(brc.db)
	uow := repositories.NewSqlUnitOfWork(brc.db)
	balanceService := services.NewBalanceService(balanceRepo, uow, nil)
	balanceHandler := handler.NewBalanceHandler(balanceService)

	brc.route.Get("/", balanceHandler.GetBalancesHandler)
	brc.route.Get("/:id", balanceHandler.GetBalanceByIDHandler)
}
