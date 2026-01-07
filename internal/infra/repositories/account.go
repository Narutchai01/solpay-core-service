package repositories

import (
	"errors"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports/repositories"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"gorm.io/gorm"
)

// Note : Ensure that GormAccountRepository implements AccountRepository
type GormAccountRepository struct {
	db *gorm.DB
}

// Note : Constructor function for GormAccountRepository
func NewGormAccountRepository(db *gorm.DB) repositories.AccountRepository {
	return &GormAccountRepository{db: db}
}

func (r *GormAccountRepository) CreateAccount(data entities.AccountEntity) (entities.AccountEntity, error) {
	// Note : Implement the logic to create an account in the database
	if err := r.db.Create(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return entities.AccountEntity{}, entities.ErrConflict
		}
		return entities.AccountEntity{}, err
	}
	return data, nil
}
