package ports

import (
	"context"

	"github.com/Narutchai01/solpay-core-service/internal/entities"
)

type AdminRepository interface {
	CreateAdmin(transaction context.Context, data *entities.AdminEntity) error
}
