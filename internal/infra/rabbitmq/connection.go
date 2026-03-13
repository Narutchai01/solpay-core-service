package rabbitmq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ holds the connection and channel to the RabbitMQ server.
type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

// NewRabbitMQ creates a new connection to the RabbitMQ server.
func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("open AMQP channel: %w", err)
	}

	log.Println("Successfully connected to RabbitMQ")

	return &RabbitMQ{
		Conn:    conn,
		Channel: ch,
	}, nil
}

// Close shuts down the channel and connection.
func (r *RabbitMQ) Close() error {
	if r.Channel != nil {
		if err := r.Channel.Close(); err != nil {
			return fmt.Errorf("close AMQP channel: %w", err)
		}
	}
	if r.Conn != nil {
		if err := r.Conn.Close(); err != nil {
			return fmt.Errorf("close AMQP connection: %w", err)
		}
	}

	log.Println("RabbitMQ connection closed")
	return nil
}
