package main

import (
	"log"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/infra/rabbitmq"
	"github.com/Narutchai01/solpay-core-service/internal/server"
)

func main() {
	cfg := config.LoadConfig()

	queueConfig := []rabbitmq.QueueConfig{
		{Name: cfg.SOLANA_WORK_QUEUE, Durable: true, AutoDelete: false},
		{Name: cfg.TRANSACTION_ORCHESTRATOR_QUEUE, Durable: true, AutoDelete: false},
		{Name: cfg.BALANCE_QUEUE, Durable: true, AutoDelete: false},
		{Name: cfg.PAYMENT_QUEUE, Durable: true, AutoDelete: false},
	}

	log.Printf("RabbitMQ URL: %s", cfg.RABBITMQ_URL)

	srv := server.New(cfg, queueConfig)
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
