package entities

import "gorm.io/gorm"

type TransactionEntity struct {
	gorm.Model
	AccountID  uint   `json:"account_id" gorm:"not null;index"`
	CategoryID uint   `json:"category_id" gorm:"not null;index"`
	Type       string `json:"type" gorm:"not null"`
	Status     string `json:"status" gorm:"not null"`
	THBAmount  int64  `json:"thb_amount" gorm:"not null;default:0"`
	USDTAmount int64  `json:"usdt_amount" gorm:"not null;default:0"`
	Fee        int64  `json:"fee" gorm:"not null;default:0"`
}
