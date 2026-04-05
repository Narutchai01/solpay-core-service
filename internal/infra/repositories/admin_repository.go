package repositories

import (
	"context"
	"errors"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/db"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"gorm.io/gorm"
)

type gormAdminRepository struct {
	db *gorm.DB
}

func NewGormAdminRepository(database *gorm.DB) ports.AdminRepository {
	return &gormAdminRepository{db: database}
}

func (r *gormAdminRepository) CreateAdmin(txCtx context.Context, data *entities.AdminEntity) error {
	tx := db.GetTx(txCtx, r.db)
	if err := tx.Create(data).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return entities.ErrConflict
		}
		return err
	}
	return nil
}
