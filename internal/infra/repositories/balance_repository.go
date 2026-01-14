package repositories

import (
	"context"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports/repositories"
	"github.com/Narutchai01/solpay-core-service/internal/db"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"gorm.io/gorm"
)

type BalanceRepository struct {
	db *gorm.DB
}

func NewGormBalanceRepository(db *gorm.DB) repositories.BalanceRepository {
	return &BalanceRepository{db: db}
}

func (r *BalanceRepository) CreateBalance(txCtx context.Context, data *entities.BalanceEntity) error {
	db := db.GetTx(txCtx, r.db)
	if err := db.Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r *BalanceRepository) GetBalances(page int, limit int) ([]entities.BalanceEntity, error) {
	var balances []entities.BalanceEntity
	offset := (page - 1) * limit

	if err := r.db.Limit(limit).Offset(offset).Find(&balances).Error; err != nil {
		return nil, err
	}
	return balances, nil
}

func (r *BalanceRepository) CountBalances() (int64, error) {
	var count int64
	if err := r.db.Model(&entities.BalanceEntity{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
