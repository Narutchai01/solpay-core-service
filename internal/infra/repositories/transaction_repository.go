package repositories

import (
	"context"
	"errors"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/db"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
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

func (r *transactionRepository) CountTransactions(accountID uint, q request.TransactionQuery) (int64, error) {
	var total int64
	db := r.db.Model(&entities.TransactionEntity{})

	if accountID > 0 {
		db = db.Where("account_id = ?", accountID)
	}
	if q.TxType != "" {
		db = db.Where("transaction_type = ?", q.TxType)
	}

	if err := db.Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *transactionRepository) GetTransactions(accountID uint, q request.TransactionQuery) ([]entities.TransactionEntity, error) {
	var transactions []entities.TransactionEntity
	db := r.db.Model(&entities.TransactionEntity{})

	if accountID > 0 {
		db = db.Where("account_id = ?", accountID)
	}
	if q.TxType != "" {
		db = db.Where("transaction_type = ?", q.TxType)
	}

	err := db.Order("created_at DESC").
		Limit(q.GetLimit()).
		Offset(q.GetOffset()).
		Preload("TransactionOnChain").
		Preload("TransactionOffChain").
		Find(&transactions).Error

	return transactions, err
}

func (r *transactionRepository) QueryTransactionSummary(txCtx context.Context, month, year int) ([]entities.TransactionSummary, error) {
	tx := db.GetTx(txCtx, r.db)

	rows, err := tx.Raw(`
		SELECT
			TO_CHAR(created_at, 'YYYY-MM-DD') AS date,
			transaction_type,
			COALESCE(
				CASE
					WHEN transaction_type = 'TOPUP'    THEN SUM(usdt_amount)
					WHEN transaction_type = 'OFFCHAIN' THEN SUM(thb_amount)
					ELSE 0
				END, 0
			) AS total_amount,
			COUNT(*) AS total_count
		FROM transaction_entities
		WHERE
			deleted_at IS NULL
			AND status = 'COMPLETED'
			AND EXTRACT(MONTH FROM created_at) = ?
			AND EXTRACT(YEAR  FROM created_at) = ?
		GROUP BY date, transaction_type
		ORDER BY date ASC
	`, month, year).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []entities.TransactionSummary
	for rows.Next() {
		var s entities.TransactionSummary
		if err := rows.Scan(&s.Date, &s.TransactionType, &s.TotalAmount, &s.TotalCount); err != nil {
			return nil, err
		}
		summaries = append(summaries, s)
	}

	return summaries, nil
}
