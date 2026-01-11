package entities

import "gorm.io/gorm"

type BalanceEntity struct {
	gorm.Model
	AccountID uint  `json:"account_id" gorm:"not null;index"`
	THBAmount int64 `json:"thb_amount" gorm:"not null;default:0"`
}
