package repositories

import (
	"github.com/Narutchai01/solpay-core-service/internal/entities"
)

type AccountRepository interface {
	CreateAccount(data entities.AccountEntity) (entities.AccountEntity, error)
}
