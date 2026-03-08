package routes

import (
	"encoding/json"
	"log"

	"github.com/Narutchai01/solpay-core-service/internal/dto/response"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/websocket"
	firberws "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SetupWebSocketRoutes(app *fiber.App, hub *websocket.Hub, db *gorm.DB) {
	// Middleware: เช็คว่าเป็น WebSocket Upgrade ไหม
	app.Use("/ws", func(c *fiber.Ctx) error {
		if firberws.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// Route: ws://localhost:PORT/ws/transaction/:tx_id
	app.Get("/ws/transaction/:tx_id", firberws.New(func(c *firberws.Conn) {

		log.Printf("New WebSocket connection for tx_id: %s", c.Params("tx_id"))
		txID := c.Params("tx_id")
		if txID == "" {
			c.Close()
			return
		}

		// สร้าง Client ใหม่
		client := &websocket.Client{
			Hub:  hub,
			Conn: c,
			Send: make(chan []byte, 256),
			TxID: txID,
		}

		// ลงทะเบียนเข้า Hub
		client.Hub.Register <- client

		// ดึงข้อมูล Transaction จาก DB แล้วส่งกลับเป็น response แรก
		go func() {
			initialPayload := fetchTransactionPayload(db, txID)
			client.Send <- initialPayload
		}()

		// แยก Process การส่งข้อความไปทำงานเงียบๆ (Goroutine)
		go client.WritePump()

		// ให้ Process หลักรอรับข้อความ (Blocking)
		client.ReadPump()
	}))
}

// fetchTransactionPayload queries the DB for the transaction and returns a JSON payload.
func fetchTransactionPayload(db *gorm.DB, txID string) []byte {
	txUUID, err := uuid.Parse(txID)
	if err != nil {
		log.Printf("Invalid tx_id format: %s, error: %v", txID, err)
		errResp, _ := json.Marshal(map[string]string{
			"error": "invalid tx_id format",
		})
		return errResp
	}

	var transaction entities.TransactionEntity
	if err := db.Preload("TransactionOnChain").Preload("TransactionOffChain").
		Where("transaction_uuid = ?", txUUID.String()).
		First(&transaction).Error; err != nil {
		log.Printf("Transaction not found for tx_id: %s, error: %v", txID, err)
		errResp, _ := json.Marshal(map[string]string{
			"error": "transaction not found",
		})
		return errResp
	}

	dto := response.FormaterTransactionDTO(&transaction)
	payload, err := json.Marshal(dto)
	if err != nil {
		log.Printf("Failed to marshal transaction for tx_id: %s, error: %v", txID, err)
		errResp, _ := json.Marshal(map[string]string{
			"error": "internal server error",
		})
		return errResp
	}

	log.Printf("Sent initial transaction data for tx_id: %s", txID)
	return payload
}
