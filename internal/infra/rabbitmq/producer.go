package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Narutchai01/solpay-core-service/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	channel *amqp.Channel
}

type ProducerInterface interface {
	SummitSolanaTx(queueName string, msg models.SolanaTxMessage) error
	PublisherTx(queueName string, body []byte) error
}

func NewProducer(channel *amqp.Channel) ProducerInterface {
	return &Producer{
		channel: channel,
	}
}

func (p *Producer) PublisherTx(queueName string, body []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ใช้ Channel ในการยิงข้อความ (Publish)
	return p.channel.PublishWithContext(ctx,
		"",        // exchange (เว้นว่างไว้คือส่งเข้า Default Exchange ตรงๆ)
		queueName, // routing key (ถ้าไม่ใส่ exchange ชื่อ routing key จะตรงกับชื่อคิว)
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent, // ⭐️ สำคัญ: บอกให้ RabbitMQ Save ข้อความลง Disk กันหาย
			Body:         body,
		},
	)
}

func (p *Producer) SummitSolanaTx(queueName string, msg models.SolanaTxMessage) error {
	q, err := p.channel.QueueDeclare(
		queueName, // ชื่อ Queue
		true,      // durable: true (ให้จำคิวไว้แม้ Server ดับ)
		false,     // autoDelete: false
		false,     // exclusive: false
		false,     // noWait: false
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	body, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("Failed to marshal message: %v", err)
	}

	p.PublisherTx(q.Name, body)
	return nil
}
