package main

import (
	"log"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/db"
	"github.com/Narutchai01/solpay-core-service/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {

		log.Fatal("Error loading .env file")
	}
	cfg := config.LoadConfig()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
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
