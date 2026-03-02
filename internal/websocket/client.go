package websocket

import (
	"log"

	"github.com/gofiber/contrib/websocket"
)

type Client struct {
	Hub  *Hub            // อ้างอิงไปยัง Hub หลัก
	Conn *websocket.Conn // Connection ของคนนี้
	Send chan []byte     // ช่องทางส่งข้อความหาคนนี้
}

// ReadPump: รอรับข้อความจาก Browser -> Server
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		// ได้ข้อความมาแล้ว ส่งเข้า Hub เพื่อกระจายต่อ (Broadcast)
		c.Hub.Broadcast <- message
	}
}

// WritePump: รอข้อความจาก Server -> Browser
func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Send
		if !ok {
			// Hub ปิด channel นี้แล้ว
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}
