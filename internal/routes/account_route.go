package routes

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/handler"
	"github.com/Narutchai01/solpay-core-service/internal/infra/repositories"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// AccountRouteConfig configures account-related routes.
type AccountRouteConfig struct {
	route fiber.Router
	db    *gorm.DB
}

// NewAccountRouteConfig creates a new AccountRouteConfig.
func NewAccountRouteConfig(route fiber.Router, db *gorm.DB) *AccountRouteConfig {
	return &AccountRouteConfig{
		route: route,
		db:    db,
	}
}

func (arc *AccountRouteConfig) Setup() {
	accountRepo := repositories.NewGormAccountRepository(arc.db)
	balanceRepo := repositories.NewGormBalanceRepository(arc.db)
	uow := repositories.NewSqlUnitOfWork(arc.db)
	accountService := services.NewAccountService(accountRepo, balanceRepo, uow)
	accountHandler := handler.NewAccountHandler(accountService)

	arc.route.Post("/", accountHandler.CreateAccountHandler)
	arc.route.Get("/", accountHandler.GetAccountsHandler)
	arc.route.Get("/:id", accountHandler.GetAccountByIDHandler)
}
