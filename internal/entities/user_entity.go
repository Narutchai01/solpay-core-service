package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	IDCard       string `json:"id_card" gorm:"not null;unique"`
	FirstName    string `json:"first_name" gorm:"not null"`
	LastName     string `json:"last_name" gorm:"not null"`
	BirthDate    string `json:"birth_date" gorm:"not null"`
	Status       string `json:"status" gorm:"not null default:'PENDING'" `
	ExpireDate   string `json:"expire_date" gorm:"not null"`
	FrontCardURL string `json:"front_card_url" gorm:"not null"`
	BackCardURL  string `json:"back_card_url" gorm:"not null"`
	KYCToken     string `json:"kyc_token" gorm:"not null"`
}

type UserStatus string

const (
	UserStatusPending  UserStatus = "PENDING"
	UserStatusApproved UserStatus = "APPROVED"
	UserStatusRejected UserStatus = "REJECTED"
)
