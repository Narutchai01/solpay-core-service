package repositories

import (
	"context"
	"errors"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports/repositories"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewGormTransactionRepository(db *gorm.DB) repositories.TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(txCtx context.Context, transactionData *entities.TransactionEntity, txOnChian *entities.TransactionOnChain, txOffChain *entities.TransactionOffChain) error {
	return r.db.WithContext(txCtx).Transaction(func(tx *gorm.DB) error {
		if transactionData == nil {
			return errors.New("transaction data is required")
		}
		if err := tx.Create(transactionData).Error; err != nil {
			return err
		}

		if txOnChian != nil {
			txOnChian.TransactionID = transactionData.TransactionUUID
			if err := tx.Create(txOnChian).Error; err != nil {
				return err
			}
		}

		if txOffChain != nil {
			txOffChain.TransactionID = transactionData.TransactionUUID
			if err := tx.Create(txOffChain).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
