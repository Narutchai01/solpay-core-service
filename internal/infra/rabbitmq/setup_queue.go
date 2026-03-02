package rabbitmq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// QueueConfig เก็บการตั้งค่าที่จำเป็นสำหรับแต่ละ Queue
type QueueConfig struct {
	Name       string
	Durable    bool
	AutoDelete bool
	// อนาคตสามารถเพิ่ม ExchangeName, RoutingKey, หรือ Arguments ตรงนี้ได้
}

// SetupQueues รับ Slice ของ QueueConfig มาวนลูปสร้าง
func SetupQueues(ch *amqp.Channel, queues []QueueConfig) error {
	for _, q := range queues {
		_, err := ch.QueueDeclare(
			q.Name,
			q.Durable,
			q.AutoDelete,
			false, // exclusive
			false, // noWait
			nil,   // arguments (เช่น ตั้งค่า Dead Letter Exchange)
		)
		if err != nil {
			// ถ้าสร้าง Queue ไหนพัง ให้ Return error กลับไปพร้อมบอกชื่อ Queue
			return fmt.Errorf("failed to declare queue '%s': %v", q.Name, err)
		}
	}

	log.Println("🐰 Successfully setup queue to RabbitMQ ✅")
	// ถ้ามี Exchange กับ Binding ก็สามารถวนลูปทำต่อด้านล่างนี้ได้เลยครับ

	return nil
}
