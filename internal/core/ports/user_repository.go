package ports

import "github.com/Narutchai01/solpay-core-service/internal/entities"

type UserRepository interface {
	CreateUser(user *entities.User) error
}
