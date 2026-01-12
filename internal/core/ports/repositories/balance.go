package repositories

import (
	"context"

	"github.com/Narutchai01/solpay-core-service/internal/entities"
)

type BalanceRepository interface {
	CreateBalance(txCtx context.Context, data *entities.BalanceEntity) error
}
