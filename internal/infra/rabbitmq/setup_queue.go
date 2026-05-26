package rabbitmq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// QueueConfig holds the configuration for a single RabbitMQ queue.
type QueueConfig struct {
	Name       string
	Durable    bool
	AutoDelete bool
}

// SetupQueues declares all queues defined in the given configuration slice.
func SetupQueues(ch *amqp.Channel, queues []QueueConfig) error {
	for _, q := range queues {
		_, err := ch.QueueDeclare(
			q.Name,
			q.Durable,
			q.AutoDelete,
			false, // exclusive
			false, // noWait
			nil,   // arguments
		)
		if err != nil {
			return fmt.Errorf("declare queue %q: %w", q.Name, err)
		}
	}

	log.Println("Successfully set up RabbitMQ queues")
	return nil
}
