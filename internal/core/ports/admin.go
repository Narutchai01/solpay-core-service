package ports

import (
	"context"

	"github.com/Narutchai01/solpay-core-service/internal/entities"
)

type AdminRepository interface {
	CreateAdmin(transaction context.Context, data *entities.AdminEntity) error
	GetAdminByUsername(ctx context.Context, username string) (*entities.AdminEntity, error)
	GetAdminByID(ctx context.Context, id string) (*entities.AdminEntity, error)
}
