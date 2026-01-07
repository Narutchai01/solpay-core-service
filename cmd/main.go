package main

import (
	"log"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/server"
)

func main() {
	cgf := config.LoadConfig()
	server := server.New(&cgf.APPPort, &cgf.TimeZone)

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
