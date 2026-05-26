package ports

import (
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
)

type UserRepository interface {
	CreateUser(user *entities.User) error
	GetUserByIDCard(idCard string) (*entities.User, error)
	GetUserByAccountID(accountID uint) (*entities.User, error)
	UpdateUser(user *entities.User) error
	GetUsers(query request.UserQuery) ([]*entities.User, error)
	CountUsers(query request.UserQuery) (int, error)
}
