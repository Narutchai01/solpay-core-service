package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	amqp "github.com/rabbitmq/amqp091-go"
)

type publisher struct {
	channel *amqp.Channel
}

// NewPublisher creates a new Publisher backed by an AMQP channel.
func NewPublisher(channel *amqp.Channel) ports.Publisher {
	return &publisher{
		channel: channel,
	}
}

func (p *publisher) Publish(queueName string, message []byte) error {
	err := p.channel.Publish(
		"",        // default exchange
		queueName, // routing key = queue name
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         message,
		},
	)
	if err != nil {
		return fmt.Errorf("publish to queue %q: %w", queueName, err)
	}
	return nil
}

func (p *publisher) PublishTransactionMessage(queue string, tx *entities.TransactionEntity, status string) {
	msg := request.TransactionMessage{
		TxID:         tx.TransactionUUID.String(),
		SourceWorker: "ORCHESTRATOR",
		Status:       status,
	}

	jsonMessage, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal balance result message: %v", err)
		return
	}

	if err := p.Publish(queue, jsonMessage); err != nil {
		log.Printf("Failed to publish balance result: %v", err)
	}
}
