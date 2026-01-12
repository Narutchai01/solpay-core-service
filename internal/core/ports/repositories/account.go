package repositories

import (
	"context"

	"github.com/Narutchai01/solpay-core-service/internal/entities"
)

type AccountRepository interface {
	CreateAccount(txCtx context.Context, data *entities.AccountEntity) error
	GetAccounts(page int, limit int) ([]entities.AccountEntity, error)
	CountAccounts() (int64, error)
	GetAccountByID(id int) (*entities.AccountEntity, error)
	GetAccountByPublicAddress(address string) (*entities.AccountEntity, error)
}
