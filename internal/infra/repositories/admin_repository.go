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

func (r *gormAdminRepository) GetAdminByUsername(ctx context.Context, username string) (*entities.AdminEntity, error) {
	var admin entities.AdminEntity

	err := r.db.WithContext(ctx).
		Where("username = ?", username).
		First(&admin).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &admin, err
}

func (r *gormAdminRepository) GetAdminByID(ctx context.Context, id string) (*entities.AdminEntity, error) {
	var admin entities.AdminEntity

	err := r.db.WithContext(ctx).First(&admin, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &admin, nil
}
