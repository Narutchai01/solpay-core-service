package ports

import (
	"context"

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
}
