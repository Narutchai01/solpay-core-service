package main

import (
	"log"
	"os"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	routes.RoutesConfig(app, db)

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT environment variable is not set")
	}

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}
