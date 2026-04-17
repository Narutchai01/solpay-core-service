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

func (r *transactionRepository) CountTransactions(query request.TransactionQuery, accountID *uint) (int64, error) {
	var total int64

	q := r.db.Model(&entities.TransactionEntity{})

	if accountID != nil {
		q = q.Where("account_id = ?", *accountID)
	}

	if query.TxType != "" {
		q = q.Where("transaction_type = ?", query.TxType)
	}

	err := q.Count(&total).Error

	return total, err
}

func (r *transactionRepository) GetTransactions(query request.TransactionQuery, accountID *uint) ([]entities.TransactionEntity, error) {
	var transactions []entities.TransactionEntity

	q := r.db.Model(&entities.TransactionEntity{})

	if accountID != nil {
		q = q.Where("account_id = ?", *accountID)
	}

	if query.TxType != "" {
		q = q.Where("transaction_type = ?", query.TxType)
	}

	err := q.Order("created_at DESC").
		Limit(query.GetLimit()).
		Offset(query.GetOffset()).
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
					WHEN transaction_type = 'ONCHAIN' THEN SUM(thb_amount)
					ELSE 0
				END, 0
			) AS total_amount,
			COALESCE(SUM(thb_amount), 0) AS total_thb_amount,
			COALESCE(SUM(usdt_amount), 0) AS total_usdt_amount,
			COALESCE(SUM(fee), 0) AS total_fee,
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
		if err := rows.Scan(&s.Date, &s.TransactionType, &s.TotalAmount, &s.TotalTHBAmount, &s.TotalUSDTAmount, &s.TotalFee, &s.TotalCount); err != nil {
			return nil, err
		}
		summaries = append(summaries, s)
	}

	return summaries, nil
}
