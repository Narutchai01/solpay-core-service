package routes

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/websocket"
	firberws "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func SetupWebSocketRoutes(app *fiber.App, hub *websocket.Hub) {
	// Middleware: เช็คว่าเป็น WebSocket Upgrade ไหม
	app.Use("/ws", func(c *fiber.Ctx) error {
		if firberws.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// Route: ws://localhost:PORT/ws
	app.Get("/ws", firberws.New(func(c *firberws.Conn) {
		// สร้าง Client ใหม่
		client := &websocket.Client{
			Hub:  hub,
			Conn: c,
			Send: make(chan []byte, 256),
		}

		// ลงทะเบียนเข้า Hub
		client.Hub.Register <- client

		// แยก Process การส่งข้อความไปทำงานเงียบๆ (Goroutine)
		go client.WritePump()

		// ให้ Process หลักรอรับข้อความ (Blocking)
		client.ReadPump()
	}))
}
