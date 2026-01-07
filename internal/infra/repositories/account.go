package repositories

import (
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"gorm.io/gorm"
)

type GormAccountRepository struct {
	db *gorm.DB
}

func NewGormAccountRepository(db *gorm.DB) *GormAccountRepository {
	return &GormAccountRepository{db: db}
}

func (r *GormAccountRepository) CreateAccount(data entities.AccountEntity) (entities.AccountEntity, error) {
	if err := r.db.Create(&data).Error; err != nil {
		return entities.AccountEntity{}, err
	}
	return data, nil
}
