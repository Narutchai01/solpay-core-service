package repositories

import (
	"context"

	"github.com/Narutchai01/solpay-core-service/internal/entities"
)

type TransactionRepository interface {
	CreateTransaction(txCtx context.Context, data *entities.TransactionEntity) error
}
