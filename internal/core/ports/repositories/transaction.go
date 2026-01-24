package repositories

import (
	"context"

	"github.com/Narutchai01/solpay-core-service/internal/entities"
)

type TransactionRepository interface {
	CreateTransaction(txCtx context.Context, transactionData *entities.TransactionEntity, txOnChian *entities.TransactionOnChain, txOffChain *entities.TransactionOffChain) error
}
