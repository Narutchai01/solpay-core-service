package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Recipient struct {
	gorm.Model
	RecipientID string `json:"recipient_id" gorm:"not null;uniqueIndex;"`
	Number      string `json:"number" gorm:"not null;"`
}

type LogPayment struct {
	gorm.Model
	TransactionUUID    uuid.UUID `json:"transaction_uuid" gorm:"not null;index;uniqueIndex;type:uuid"`
	OmiseTransactionID string    `json:"omise_transaction_id" gorm:"not null;index"`
}
