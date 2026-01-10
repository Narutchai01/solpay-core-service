package repositories

import (
	"github.com/Narutchai01/solpay-core-service/internal/entities"
)

type BalanceRepository interface {
	CreateBalance(data *entities.BalanceEntity) error
}
