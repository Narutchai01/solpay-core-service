package main

import (
	"log"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/db"
	"github.com/Narutchai01/solpay-core-service/internal/models/response"
	"github.com/Narutchai01/solpay-core-service/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	cfg := config.LoadConfig()
	// ประกาศไว้ข้างนอก Handler หรือใน main setup
	app := fiber.New()

	app.Use(logger.New(logger.Config{
		// Format:   "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeZone: cfg.TimeZone,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(response.NewResponseModel(fiber.StatusOK, "Server is running", nil, nil))
	})

	db, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	routes.RoutesConfig(app, db)

	if err := app.Listen(":" + cfg.APPPort); err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}
