package main

import (
	"fmt"
	"log"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/infra/rabbitmq"
	"github.com/Narutchai01/solpay-core-service/internal/server"
)

func main() {
	cgf := config.LoadConfig()

	queueConfig := []rabbitmq.QueueConfig{
		{Name: cgf.SOLANA_WORK_QUEUE, Durable: true, AutoDelete: false},
		{Name: cgf.TRANSACTION_ORCHESTRATOR_QUEUE, Durable: true, AutoDelete: false},
		{Name: cgf.BALANCE_QUEUE, Durable: true, AutoDelete: false},
	}

	println(fmt.Sprintf("RabbitMQ URL: %s", cgf.RABBITMQ_URL))
	server := server.New(cgf.APPPort, cgf.TimeZone, cgf.RABBITMQ_URL, queueConfig)

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
