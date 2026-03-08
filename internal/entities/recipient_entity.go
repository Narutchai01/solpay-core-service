package entities

import "gorm.io/gorm"

type Recipient struct {
	gorm.Model
	RecipientID string `json:"recipient_id" gorm:"not null;uniqueIndex;"`
	Number      string `json:"number" gorm:"not null;"`
}
