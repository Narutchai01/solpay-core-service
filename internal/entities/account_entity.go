package entities

import "gorm.io/gorm"

type AccountEntity struct {
	gorm.Model
	PublicAddress string  `json:"public_address" gorm:"not null;unique"`
	KycToken      *string `json:"kyc_token" gorm:"unique;default:null"`
	User          *User   `json:"user,omitempty" gorm:"foreignKey:AccountID;references:ID"`
}
