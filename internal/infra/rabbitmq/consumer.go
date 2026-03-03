package rabbitmq

import (
	"context"
	"log"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	ch                 *amqp.Channel
	transactionService services.TransactionService
}

func NewConsumer(ch *amqp.Channel, transactionService services.TransactionService) ports.Consumer {
	return &Consumer{
		ch:                 ch,
		transactionService: transactionService,
	}
}

func (c *Consumer) TransactionOrchestrator() error {
	cfg := config.LoadConfig().TRANSACTION_ORCHESTRATOR_QUEUE
	msgs, err := c.ch.Consume(
		cfg,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for msg := range msgs {
		log.Printf("Received message from queue: %s", string(msg.Body))
		err := c.transactionService.HandleTransactionProcess(context.Background(), msg.Body)
		if err != nil {
			log.Printf("Failed to handle transaction process: %v", err)
		}
	}

	return nil
}
