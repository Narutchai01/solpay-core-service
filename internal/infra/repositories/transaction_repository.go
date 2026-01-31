package repositories

import (
	"context"
	"errors"

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

func (r *TransactionRepository) CreateTransactionOnChain(txCtx context.Context, data *entities.TransactionOnChain) error {
	db := db.GetTx(txCtx, r.db)
	if err := db.Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r *TransactionRepository) CreateTransactionOffChain(txCtx context.Context, data *entities.TransactionOffChain) error {
	db := db.GetTx(txCtx, r.db)
	if err := db.Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r *TransactionRepository) UpdateTransactionStatus(txCtx context.Context, transactionUUID string, status string) error {
	db := db.GetTx(txCtx, r.db)
	if err := db.Model(&entities.TransactionEntity{}).Where("transaction_uuid = ?", transactionUUID).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

func (r *TransactionRepository) GetTransactionByAccountID(accountID int) ([]entities.TransactionEntity, error) {
	var transactions []entities.TransactionEntity
	if err := r.db.Where("account_id = ?", accountID).Find(&transactions).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entities.TransactionEntity{}, entities.ErrNotFound
		}
	}
	return transactions, nil
}

func (r *TransactionRepository) GetTransactionByID(transactionID int) (*entities.TransactionEntity, error) {
	var transaction entities.TransactionEntity
	if err := r.db.First(&transaction, transactionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &entities.TransactionEntity{}, entities.ErrNotFound
		}
		return &entities.TransactionEntity{}, entities.ErrInternal
	}
	return &transaction, nil
}
