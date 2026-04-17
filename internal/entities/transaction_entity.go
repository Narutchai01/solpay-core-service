package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionEntity struct {
	TransactionUUID uuid.UUID `json:"transaction_uuid" gorm:"primaryKey;not null;uniqueIndex;type:uuid"`
	gorm.Model
	AccountID           uint                 `json:"account_id" gorm:"not null;index"`
	CategoryID          string               `json:"category_id" gorm:"default:null;"`
	TransactionType     string               `json:"transaction_type" gorm:"not null;index"`
	Status              string               `json:"status" gorm:"not null;default:'pending'"`
	THBAmount           float64              `json:"thb_amount" gorm:"not null;default:0"`
	USDTAmount          float64              `json:"usdt_amount" gorm:"not null;default:0"`
	Fee                 float64              `json:"fee" gorm:"not null;default:0"`
	TransactionOnChain  *TransactionOnChain  `json:"transaction_on_chain,omitempty" gorm:"foreignKey:TransactionID;references:TransactionUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TransactionOffChain *TransactionOffChain `json:"transaction_off_chain,omitempty" gorm:"foreignKey:TransactionID;references:TransactionUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	LogPayment *LogPayment `json:"log_payment,omitempty" gorm:"foreignKey:TransactionUUID;references:TransactionUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type TransactionOnChain struct {
	gorm.Model
	TransactionID uuid.UUID          `json:"transaction_id" gorm:"type:uuid;not null;index"`
	Transaction   *TransactionEntity `json:"transaction,omitempty" gorm:"foreignKey:TransactionID;references:TransactionUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TxHash        string             `json:"tx_hash" gorm:"not null;uniqueIndex"`
}

type TransactionOffChain struct {
	gorm.Model
	TransactionID uuid.UUID          `json:"transaction_id" gorm:"type:uuid;not null;index"`
	Transaction   *TransactionEntity `json:"transaction,omitempty" gorm:"foreignKey:TransactionID;references:TransactionUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PromptPayID   string             `json:"prompt_pay_id" gorm:"not null;"`
	SlipURL       *string            `json:"slip_url" default:"null" `
}

type TransactionStatus string

const (
	StatusPending            TransactionStatus = "PENDING"              // เริ่มต้น
	StatusSolanaSubmitted    TransactionStatus = "SOLANA_SUBMITTED"     // ส่งไป Solana Worker แล้ว
	StatusSolanaSuccess      TransactionStatus = "BLOCKCHAIN_COMPLETED" // Solana ตัดเงินสำเร็จ
	StatusSolanaFailed       TransactionStatus = "SOLANA_FAILED"        // Solana ตัดเงินไม่ผ่าน
	StatusBalanceUpdating    TransactionStatus = "BALANCE_UPDATING"     // กำลังอัปเดตยอดเงิน
	StatusBalanceUpdated     TransactionStatus = "BALANCE_UPDATED"      // อัปเดตยอดเงินเสร็จแล้ว
	StatusBalanceWithdrawing TransactionStatus = "BALANCE_WITHDRAWING"  // กำลังหักยอดเงิน (สำหรับถอนเงิน)
	StatusBalanceWithdrawn   TransactionStatus = "BALANCE_WITHDRAWN"    // ยอดเงินถูกหักแล้ว (สำหรับถอนเงิน)
	StatusRefunding          TransactionStatus = "REFUNDING"            // กำลังคืนเงิน
	StatusBalanceFailed      TransactionStatus = "BALANCE_FAILED"       // อัปเดตยอดเงินไม่ผ่าน
	StatusCompleted          TransactionStatus = "COMPLETED"            // ทุกอย่างสมบูรณ์
	StatusRefunded           TransactionStatus = "REFUNDED"             // เกิดข้อผิดพลาดและคืนเงินแล้ว
	StatusPaymentSuccess     TransactionStatus = "PAYMENT_SUCCESS"
	StatusFailed             TransactionStatus = "FAILED"
	StatusPaymentFailed      TransactionStatus = "PAYMENT_FAILED"
)

type TransactionType string

const (
	TOPUP    TransactionType = "TOPUP"
	ONCHAIN  TransactionType = "ONCHAIN"
	OFFCHAIN TransactionType = "OFFCHAIN"
)

func (t *TransactionEntity) BeforeCreate(tx *gorm.DB) (err error) {
	if t.TransactionUUID == uuid.Nil {
		t.TransactionUUID = uuid.New()
	}
	return
}

type TransactionSummary struct {
	Date            string  `json:"date"`
	TransactionType string  `json:"transaction_type"`
	TotalAmount     float64 `json:"total_amount"`
	TotalTHBAmount  float64 `json:"total_thb_amount"`
	TotalUSDTAmount float64 `json:"total_usdt_amount"`
	TotalFee        float64 `json:"total_fee"`
	TotalCount      int     `json:"total_count"`
}
