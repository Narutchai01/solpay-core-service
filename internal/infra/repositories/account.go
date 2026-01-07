package repositories

import (
	"errors"

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
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return entities.AccountEntity{}, entities.ErrConflict
		}
		return entities.AccountEntity{}, err
	}
	return data, nil
}
