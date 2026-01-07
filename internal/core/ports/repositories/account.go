package repositories

import (
	"github.com/Narutchai01/solpay-core-service/internal/entities"
)

type AccountRepository interface {
	CreateAccount(data *entities.AccountEntity) error
	GetAccounts(page int, limit int) ([]entities.AccountEntity, error)
	CountAccounts() (int64, error)
}
