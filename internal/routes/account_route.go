package routes

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/handler"
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/infra/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AccountRouteConfig struct {
	route    fiber.Router
	db       *gorm.DB
	validate *validator.Validate
}

func NewAccountRouteConfig(route fiber.Router, db *gorm.DB, validate *validator.Validate) *AccountRouteConfig {
	return &AccountRouteConfig{
		route:    route,
		db:       db,
		validate: validate,
	}
}

func (arc *AccountRouteConfig) Setup() {
	accountRepository := repositories.NewGormAccountRepository(arc.db)
	accountService := services.NewAccountService(accountRepository)
	accountHandler := handler.NewAccountHandler(accountService)

	arc.route.Post("/", accountHandler.CreateAccountHandler)
	arc.route.Get("/", accountHandler.GetAccountsHandler)
	arc.route.Get("/:id", accountHandler.GetAccountByIDHandler)
	arc.route.Get("/address/:publicAddress", accountHandler.GetAccountByPublicAddressHandler)
}
