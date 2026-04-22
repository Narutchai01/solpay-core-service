package ports

import "github.com/Narutchai01/solpay-core-service/internal/entities"

type UserRepository interface {
	CreateUser(user *entities.User) error
	GetUserByIDCard(idCard string) (*entities.User, error)
	UpdateUser(user *entities.User) error
}
