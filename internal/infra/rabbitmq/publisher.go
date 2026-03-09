package rabbitmq

import (
	"fmt"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
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
