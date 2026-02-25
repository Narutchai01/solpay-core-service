package rabbitmq

import (
	"github.com/Narutchai01/solpay-core-service/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	channel *amqp.Channel
}

type ProducerInterface interface {
	SummitSolanaTx(queueName string, msg models.SolanaTxMessage) error
}

func NewProducer(channel *amqp.Channel) ProducerInterface {
	return &Producer{
		channel: channel,
	}
}

func (p *Producer) SummitSolanaTx(queueName string, msg models.SolanaTxMessage) error {

	println("send msg")

	p.channel.Publish(
		"",        // exchange
		queueName, // routing key (queue)
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(msg.Base64Tx),
		},
	)
	return nil
}
