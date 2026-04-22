package repositories

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"gorm.io/gorm"
)

func NewGormUserRepository(db *gorm.DB) ports.UserRepository {
	return &gormUserRepository{db: db}
}

type gormUserRepository struct {
	db *gorm.DB
}

func (r *gormUserRepository) CreateUser(user *entities.User) error {
	return r.db.Create(user).Error
}
