package rabbitmq

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	channel *amqp.Channel
}

func NewPublisher(channel *amqp.Channel) ports.Publisher {
	return &Publisher{
		channel: channel,
	}
}

func (p *Publisher) Publish(queueName string, message []byte) error {
	return p.channel.Publish(
		"",        // exchange (เว้นว่างไว้คือส่งเข้า Default Exchange ตรงๆ)
		queueName, // routing key (ถ้าไม่ใส่ exchange ชื่อ routing key จะตรงกับชื่อคิว)
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent, // ⭐️ สำคัญ: บอกให้ RabbitMQ Save ข้อความลง Disk กันหาย
			Body:         message,
		},
	)
}
