package rabbitmq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ struct เก็บ Connection และ Channel เพื่อนำไปใช้ต่อ
type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

// NewRabbitMQ ทำหน้าที่สร้างการเชื่อมต่อไปยัง RabbitMQ server
func NewRabbitMQ(url string) (*RabbitMQ, error) {
	// 1. สร้าง Connection
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// 2. สร้าง Channel (เปรียบเสมือนท่อส่งข้อมูลภายใน Connection)
	ch, err := conn.Channel()
	if err != nil {
		// ปิด Connection ถ้าสร้าง Channel ไม่สำเร็จ
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	log.Println("🐰 Successfully connected to RabbitMQ ✅")

	return &RabbitMQ{
		Conn:    conn,
		Channel: ch,
	}, nil
}

// Close ใช้สำหรับปิด Channel และ Connection เมื่อเลิกใช้งาน
func (r *RabbitMQ) Close() error {
	if r.Channel != nil {
		if err := r.Channel.Close(); err != nil {
			return fmt.Errorf("failed to close channel: %w", err)
		}
	}
	if r.Conn != nil {
		if err := r.Conn.Close(); err != nil {
			return fmt.Errorf("failed to close connection: %w", err)
		}
	}

	log.Println("🛑 RabbitMQ connection closed")
	return nil
}
