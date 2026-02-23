package server

import (
	"log"

	"github.com/Narutchai01/solpay-core-service/internal/core/websocket"
	"github.com/Narutchai01/solpay-core-service/internal/db"
	"github.com/Narutchai01/solpay-core-service/internal/infra/rabbitmq"
	"github.com/Narutchai01/solpay-core-service/internal/routes"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Server struct {
	App          *fiber.App
	Port         string
	TimeZone     string
	RABBITMQ_URL string
}

func New(port string, timeZone string, rabbitMQURL string) *Server {
	app := fiber.New(fiber.Config{
		AppName: "Solpay core service",
	})

	return &Server{
		App:          app,
		Port:         port,
		TimeZone:     timeZone,
		RABBITMQ_URL: rabbitMQURL,
	}
}

func (s *Server) Start() error {

	s.App.Use(logger.New(logger.Config{
		// Format:   "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeZone: s.TimeZone,
	}))

	s.App.Get("/", func(c *fiber.Ctx) error {
		return utils.HandleResponse(c, nil, nil, "Server is running")
	})

	hub := websocket.NewHub()
	go hub.Run()
	routes.SetupWebSocketRoutes(s.App, hub)

	db, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	_, err = rabbitmq.NewRabbitMQ(s.RABBITMQ_URL)
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %v", err)
	}

	routes.RoutesConfig(s.App, db)

	return s.App.Listen(":" + s.Port)
}
