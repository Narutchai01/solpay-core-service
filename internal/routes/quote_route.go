package routes

import (
	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/handler"
	"github.com/Narutchai01/solpay-core-service/internal/infra/repositories"
	"github.com/Narutchai01/solpay-core-service/internal/infra/solana"
	"github.com/Narutchai01/solpay-core-service/internal/middlewares"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type QuoteRouteConfig struct {
	route fiber.Router
	db    *gorm.DB
}

func NewQuoteRouteConfig(route fiber.Router, db *gorm.DB) *QuoteRouteConfig {
	return &QuoteRouteConfig{
		route: route,
		db:    db,
	}
}

func (qrc *QuoteRouteConfig) Setup() {
	cfg := config.LoadConfig()
	// สร้าง Quote Repository, Service, และ Handler
	solanaRepo := solana.NewSolanaTransactionRepository(cfg.RPC_URL, cfg.MINT_TOKEN_ADDRESS)
	accountRepo := repositories.NewGormAccountRepository(qrc.db)
	quoteRepo := repositories.NewGormQuoteRepository(qrc.db)
	quoteService := services.NewQuoteService(quoteRepo, solanaRepo, accountRepo) // ผ่าน solanaRepo เป็น nil ชั่วคราว
	quoteHandler := handler.NewQuoteHandler(quoteService)

	// กำหนดเส้นทางสำหรับ Quote
	qrc.route.Post("/", middlewares.AuthRequired(), quoteHandler.CreateQuoteHandler)
	qrc.route.Patch("/:id/confirm", middlewares.AuthRequired(), quoteHandler.ConFirmQuoteHandler)

	qrc.route.Get("/:id", quoteHandler.GetQuoteByIDHandler)
}
