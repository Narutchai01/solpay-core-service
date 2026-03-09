package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	amqp "github.com/rabbitmq/amqp091-go"
)

type consumer struct {
	ch                 *amqp.Channel
	transactionService services.TransactionService
	balanceService     services.BalanceService
	paymentService     services.PaymentService
	publisher          ports.Publisher
}

// NewConsumer creates a new Consumer that processes messages from RabbitMQ queues.
func NewConsumer(ch *amqp.Channel, transactionService services.TransactionService, balanceService services.BalanceService, paymentService services.PaymentService, publisher ports.Publisher) ports.Consumer {
	return &consumer{
		ch:                 ch,
		transactionService: transactionService,
		balanceService:     balanceService,
		paymentService:     paymentService,
		publisher:          publisher,
	}
}

func (c *consumer) TransactionOrchestrator() error {
	cfg := config.LoadConfig()
	msgs, err := c.ch.Consume(
		cfg.TRANSACTION_ORCHESTRATOR_QUEUE,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("consume transaction orchestrator queue: %w", err)
	}

	for msg := range msgs {
		log.Printf("Received message from transaction orchestrator queue: %s", string(msg.Body))
		if err := c.transactionService.HandleTransactionUpdate(context.Background(), msg.Body); err != nil {
			log.Printf("Failed to handle transaction update: %v", err)
		}
	}

	return nil
}

func (c *consumer) BalanceConsumer() error {
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
		return fmt.Errorf("consume balance queue: %w", err)
	}

	for msg := range msgs {
		log.Printf("Received message from balance queue: %s", string(msg.Body))
		if err := c.balanceService.UpdateBalance(msg.Body); err != nil {
			log.Printf("Failed to handle balance update: %v", err)
		}
	}

	return nil
}

func (c *consumer) PaymentConsumer() error {
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
		return fmt.Errorf("consume payment queue: %w", err)
	}

	for msg := range msgs {
		log.Printf("Received message from payment queue: %s", string(msg.Body))
		data, err := c.paymentService.ProcessPayment(context.Background(), msg.Body)
		if err != nil {
			log.Printf("Failed to handle payment: %v", err)

			payload := request.TransactionMessage{
				TxID:         "",
				SourceWorker: "PAYMENT",
				Status:       string(entities.StatusPaymentFailed),
			}

			failBody, marshalErr := json.Marshal(payload)
			if marshalErr != nil {
				log.Printf("Failed to marshal failure payload: %v", marshalErr)
				continue
			}

			if pubErr := c.publisher.Publish(cfg.TRANSACTION_ORCHESTRATOR_QUEUE, failBody); pubErr != nil {
				log.Printf("Failed to publish failure message: %v", pubErr)
			}
			continue
		}

		if err := c.publisher.Publish(cfg.TRANSACTION_ORCHESTRATOR_QUEUE, data); err != nil {
			log.Printf("Failed to publish payment result: %v", err)
		}
	}

	return nil
}
