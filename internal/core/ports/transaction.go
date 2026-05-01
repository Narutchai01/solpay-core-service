package ports

import (
	"context"

	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/google/uuid"
)

type TransactionRepository interface {
	CreateTransaction(txCtx context.Context, data *entities.TransactionEntity) error
	CreateTransactionOnChain(txCtx context.Context, data *entities.TransactionOnChain) error
	CreateTransactionOffChain(txCtx context.Context, data *entities.TransactionOffChain) error
	UpdateTransactionStatus(txCtx context.Context, transactionUUID string, status string) error
	GetTransactionByID(id int) (*entities.TransactionEntity, error)
	GetTransactionByAccountID(accountID int) ([]entities.TransactionEntity, error)
	GetTransactionByUUID(txUUID uuid.UUID) (*entities.TransactionEntity, error)
	CountTransactions(query request.TransactionQuery, accountID *uint) (int64, error)
	GetTransactions(query request.TransactionQuery, accountID *uint) ([]entities.TransactionEntity, error)
	QueryTransactionSummary(txCtx context.Context, month, year int) ([]entities.TransactionSummary, error)
	GetSpendingSummary(ctx context.Context, accountID uint, month, year int) ([]entities.SpendingSummary, error)
	GetMonthlySpendingSummary(ctx context.Context, accountID uint, limit int) ([]entities.MonthlySpending, error)
}
