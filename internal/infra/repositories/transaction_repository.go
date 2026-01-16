package repositories

import (
	"context"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports/repositories"
	"github.com/Narutchai01/solpay-core-service/internal/db"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewGormTransactionRepository(db *gorm.DB) repositories.TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(txCtx context.Context, data *entities.TransactionEntity) error {
	db := db.GetTx(txCtx, r.db)
	if err := db.Create(&data).Error; err != nil {
		return err
	}
	return nil
}
