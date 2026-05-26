package entities

import (
	"gorm.io/gorm"
)

type AdminEntity struct {
	gorm.Model
	Username string `json:"username" gorm:"type:text;not null;unique"`
	Password string `json:"password" gorm:"type:text;not null"`
}
