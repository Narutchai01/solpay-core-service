package main

import (
	"fmt"
	"log"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/server"
)

func main() {
	cgf := config.LoadConfig()

	println(fmt.Sprintf("RabbitMQ URL: %s", cgf.RABBITMQ_URL))
	server := server.New(cgf.APPPort, cgf.TimeZone, cgf.RABBITMQ_URL)

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
