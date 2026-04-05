package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminEntity struct {
	gorm.Model
	ID       uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username string    `json:"username" gorm:"type:text;not null;unique"`
	Password string    `json:"password" gorm:"type:text;not null"`
}
