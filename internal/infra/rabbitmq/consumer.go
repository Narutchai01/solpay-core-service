package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	ch                 *amqp.Channel
	transactionService services.TransactionService
	balanceService     services.BalanceService
	paymentService     services.PaymentService
	publisher          ports.Publisher
}

func NewConsumer(ch *amqp.Channel, transactionService services.TransactionService, balanceService services.BalanceService, paymentService services.PaymentService, publisher ports.Publisher) ports.Consumer {
	return &Consumer{
		ch:                 ch,
		transactionService: transactionService,
		balanceService:     balanceService,
		paymentService:     paymentService,
		publisher:          publisher,
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

func (c *Consumer) PaymentConsumer() error {
	cfg := config.LoadConfig()
	msgs, err := c.ch.Consume(
		cfg.PAYMENT_QUEUE,
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
		log.Printf("Received message from payment queue: %s", string(msg.Body))
		data, err := c.paymentService.ProcessPayment(context.Background(), msg.Body)
		if err != nil {
			log.Printf("Failed to handle payment: %v", err)

			payload := request.TransactionMessage{
				TxID:         "",
				SourceWorker: "PAYMENT",
				Status:       "PAYMENT_FAILED",
			}
			msg.Body, _ = json.Marshal(payload)

			c.publisher.Publish(cfg.TRANSACTION_ORCHESTRATOR_QUEUE, msg.Body) // Re-queue the message for retry
			continue
		}

		err = c.publisher.Publish(cfg.TRANSACTION_ORCHESTRATOR_QUEUE, data)
		if err != nil {
			log.Printf("Failed to publish payment result: %v", err)
		}
	}
	return nil
}
