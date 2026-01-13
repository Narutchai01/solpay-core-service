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
