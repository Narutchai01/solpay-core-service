package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Quote struct {
	ID           string    `gorm:"primaryKey;type:varchar(50)"` // เช่น q_1234...
	AccountID    int64     `gorm:"index"`
	Type         string    `gorm:"type:varchar(20)"`                   // เช่น "TOPUP_CRYPTO"
	THBAmount    int64     `gorm:"column:thb_amount"`                  // เก็บเป็นสตางค์ เช่น 100000 (1,000 บาท)
	USDTAmount   float64   `gorm:"column:usdt_amount"`                 // เช่น 30.88
	ExchangeRate float64   `gorm:"column:exchange_rate"`               // เช่น 32.39
	ExpiresAt    time.Time `gorm:"index"`                              // เวลาหมดอายุ (สำคัญมาก)
	Status       string    `gorm:"type:varchar(20);default:'PENDING'"` // PENDING, USED, EXPIRED
	CreatedAt    time.Time
	Fee          float64 `gorm:"column:fee"` // ค่าธรรมเนียม (ถ้ามี)
}

// สร้าง ID อัตโนมัติก่อน Save
func (q *Quote) BeforeCreate(tx *gorm.DB) (err error) {
	q.ID = "q_" + uuid.New().String()
	return
}
