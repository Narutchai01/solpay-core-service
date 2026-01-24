package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionEntity struct {
	gorm.Model
	TransactionUUID     uuid.UUID            `json:"transaction_uuid" gorm:"not null;uniqueIndex;type:uuid"`
	AccountID           uint                 `json:"account_id" gorm:"not null;index"`
	TransactionType     string               `json:"category_id" gorm:"not null;"`
	Status              string               `json:"status" gorm:"not null;default:'pending'"`
	THBAmount           float64              `json:"thb_amount" gorm:"not null;default:0"`
	USDTAmount          float64              `json:"usdt_amount" gorm:"not null;default:0"`
	Fee                 float64              `json:"fee" gorm:"not null;default:0"`
	TransactionOnChain  *TransactionOnChain  `json:"transaction_on_chain,omitempty" gorm:"foreignKey:TransactionID;references:TransactionUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TransactionOffChain *TransactionOffChain `json:"transaction_off_chain,omitempty" gorm:"foreignKey:TransactionID;references:TransactionUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type TransactionOnChain struct {
	gorm.Model
	TransactionID uuid.UUID          `json:"transaction_id" gorm:"type:uuid;not null;index"`
	Transaction   *TransactionEntity `json:"transaction,omitempty" gorm:"foreignKey:TransactionID;references:TransactionUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TxHash        string             `json:"tx_hash" gorm:"not null;uniqueIndex"`
	FromAddress   string             `json:"from_address" gorm:"not null"`
}

type TransactionOffChain struct {
	gorm.Model
	TransactionID uuid.UUID          `json:"transaction_id" gorm:"type:uuid;not null;index"`
	Transaction   *TransactionEntity `json:"transaction,omitempty" gorm:"foreignKey:TransactionID;references:TransactionUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PropmtPayID   string             `json:"prompt_pay_id" gorm:"not null;"`
	SlipURL       *string            `json:"slip_url" default:"null" `
}
