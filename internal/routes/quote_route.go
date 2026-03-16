package routes

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/handler"
	"github.com/Narutchai01/solpay-core-service/internal/infra/repositories"
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
	// สร้าง Quote Repository, Service, และ Handler
	quoteRepo := repositories.NewGormQuoteRepository(qrc.db)
	quoteService := services.NewQuoteService(quoteRepo)
	quoteHandler := handler.NewQuoteHandler(quoteService)

	// กำหนดเส้นทางสำหรับ Quote
	qrc.route.Post("/", quoteHandler.CreateQuoteHandler)
	qrc.route.Get("/:id", quoteHandler.GetQuoteByIDHandler)
}
