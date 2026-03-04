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
	balanceService     services.BalanceService
}

func NewConsumer(ch *amqp.Channel, transactionService services.TransactionService, balanceService services.BalanceService) ports.Consumer {
	return &Consumer{
		ch:                 ch,
		transactionService: transactionService,
		balanceService:     balanceService,
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
		err := c.transactionService.HandleTransactionUpdate(context.Background(), msg.Body)
		if err != nil {
			log.Printf("Failed to handle transaction update: %v", err)
		}
	}

	return nil
}

func (c *Consumer) BalanceConsumer() error {
	cfg := config.LoadConfig()
	msgs, err := c.ch.Consume(
		cfg.BALANCE_QUEUE,
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
		log.Printf("Received message from balance queue: %s", string(msg.Body))
		err := c.balanceService.UpdateBalance(msg.Body)
		if err != nil {
			log.Printf("Failed to handle balance update: %v", err)
		}
	}

	return nil
}
