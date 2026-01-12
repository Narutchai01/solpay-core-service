package entities

import "gorm.io/gorm"

type BalanceEntity struct {
	gorm.Model
	AccountID  uint          `json:"account_id" gorm:"not null;index"`
	THBAmount  int64         `json:"thb_amount" gorm:"not null;default:0"`
	USDTAmount int64         `json:"usdt_amount" gorm:"not null;default:0"`
	Account    AccountEntity `json:"account" gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
