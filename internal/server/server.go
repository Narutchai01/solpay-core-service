package server

import (
	"log"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/db"
	"github.com/Narutchai01/solpay-core-service/internal/infra/rabbitmq"
	"github.com/Narutchai01/solpay-core-service/internal/routes"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/Narutchai01/solpay-core-service/internal/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/omise/omise-go"
)

// Server represents the HTTP server and its dependencies.
type Server struct {
	app         *fiber.App
	cfg         *config.Config
	queueConfig []rabbitmq.QueueConfig
}

// New creates a new Server with the given configuration.
func New(cfg *config.Config, queueConfig []rabbitmq.QueueConfig) *Server {
	app := fiber.New(fiber.Config{
		AppName: "Solpay core service",
	})

	return &Server{
		app:         app,
		cfg:         cfg,
		queueConfig: queueConfig,
	}
}

// Start initialises all dependencies and starts the HTTP server.
func (s *Server) Start() error {
	s.app.Use(logger.New(logger.Config{
		TimeZone: s.cfg.TimeZone,
	}))

	s.app.Use(cors.New(cors.Config{
		AllowOrigins:     s.cfg.ALLOW_ORIGIN,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	s.app.Get("/", func(c *fiber.Ctx) error {
		return utils.HandleResponse(c, nil, nil, "Server is running")
	})

	hub := websocket.NewHub()
	go hub.Run()

	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	routes.SetupWebSocketRoutes(s.app, hub, dbConn)

	mq, err := rabbitmq.NewRabbitMQ(s.cfg.RABBITMQ_URL)
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %v", err)
	}
	defer mq.Close()

	if err := rabbitmq.SetupQueues(mq.Channel, s.queueConfig); err != nil {
		log.Fatalf("Error setting up RabbitMQ queues: %v", err)
	}

	omiseClient, err := omise.NewClient(s.cfg.OMISE_KEY, s.cfg.OMISE_SECRET)
	if err != nil {
		log.Fatalf("Error creating Omise client: %v", err)
	}

	consumerSetup := rabbitmq.NewConsumerSetup(mq.Channel, dbConn, hub, omiseClient)
	consumerSetup.Setup()

	routes.RoutesConfig(s.app, dbConn, mq.Channel, s.cfg)

	return s.app.Listen(":" + s.cfg.APPPort)
}
