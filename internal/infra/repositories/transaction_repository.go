package repositories

import (
	"context"
	"errors"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/db"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

// NewGormTransactionRepository creates a new GORM-backed TransactionRepository.
func NewGormTransactionRepository(database *gorm.DB) ports.TransactionRepository {
	return &transactionRepository{db: database}
}

func (r *transactionRepository) CreateTransaction(txCtx context.Context, data *entities.TransactionEntity) error {
	tx := db.GetTx(txCtx, r.db)
	return tx.Create(data).Error
}

func (r *transactionRepository) CreateTransactionOnChain(txCtx context.Context, data *entities.TransactionOnChain) error {
	tx := db.GetTx(txCtx, r.db)
	return tx.Create(data).Error
}

func (r *transactionRepository) CreateTransactionOffChain(txCtx context.Context, data *entities.TransactionOffChain) error {
	tx := db.GetTx(txCtx, r.db)
	return tx.Create(data).Error
}

func (r *transactionRepository) UpdateTransactionStatus(txCtx context.Context, transactionUUID string, status string) error {
	tx := db.GetTx(txCtx, r.db)
	return tx.Model(&entities.TransactionEntity{}).
		Where("transaction_uuid = ?", transactionUUID).
		Update("status", status).Error
}

func (r *transactionRepository) GetTransactionByAccountID(accountID int) ([]entities.TransactionEntity, error) {
	var transactions []entities.TransactionEntity
	if err := r.db.Where("account_id = ?", accountID).Find(&transactions).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrNotFound
		}
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) GetTransactionByID(transactionID int) (*entities.TransactionEntity, error) {
	var transaction entities.TransactionEntity
	if err := r.db.First(&transaction, transactionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrNotFound
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) GetTransactionByUUID(txUUID uuid.UUID) (*entities.TransactionEntity, error) {
	var transaction entities.TransactionEntity
	err := r.db.
		Preload("TransactionOnChain").
		Preload("TransactionOffChain").
		Where("transaction_uuid = ?", txUUID.String()).
		First(&transaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrNotFound
		}
		return nil, err
	}
	return &transaction, nil
}
