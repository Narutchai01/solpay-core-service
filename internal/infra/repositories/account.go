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

func (r *GormAccountRepository) CreateAccount(data *entities.AccountEntity) error {
	// Note : Implement the logic to create an account in the database
	if err := r.db.Create(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return entities.ErrConflict
		}
		return err
	}
	return nil
}

func (r *GormAccountRepository) GetAccounts(page int, limit int) ([]entities.AccountEntity, error) {
	var accounts []entities.AccountEntity
	offset := (page - 1) * limit

	if err := r.db.Limit(limit).Offset(offset).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *GormAccountRepository) CountAccounts() (int64, error) {
	var count int64
	if err := r.db.Model(&entities.AccountEntity{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *GormAccountRepository) GetAccountByID(accountID int) (entities.AccountEntity, error) {
	var account entities.AccountEntity
	if err := r.db.First(&account, accountID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.AccountEntity{}, entities.ErrNotFound
		}
		return entities.AccountEntity{}, entities.ErrInternal
	}
	return account, nil
}
