package repositories

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/ports/repositories"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"gorm.io/gorm"
)

type BalanceRepository struct {
	db *gorm.DB
}

func NewBalanceRepository(db *gorm.DB) repositories.BalanceRepository {
	return &BalanceRepository{db: db}
}

func (r *BalanceRepository) CreateBalance(data *entities.BalanceEntity) error {
	if err := r.db.Create(&data).Error; err != nil {
		return err
	}
	return nil
}
