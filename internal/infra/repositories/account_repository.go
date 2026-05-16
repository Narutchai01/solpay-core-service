package repositories

import (
	"context"
	"errors"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/db"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"gorm.io/gorm"
)

type gormAccountRepository struct {
	db *gorm.DB
}

// NewGormAccountRepository creates a new GORM-backed AccountRepository.
func NewGormAccountRepository(database *gorm.DB) ports.AccountRepository {
	return &gormAccountRepository{db: database}
}

func (r *gormAccountRepository) CreateAccount(txCtx context.Context, data *entities.AccountEntity) error {
	tx := db.GetTx(txCtx, r.db)
	if err := tx.Create(data).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return entities.ErrConflict
		}
		return err
	}
	return nil
}

func (r *gormAccountRepository) GetAccounts(page, limit int) ([]entities.AccountEntity, error) {
	var accounts []entities.AccountEntity
	offset := (page - 1) * limit

	if err := r.db.Limit(limit).Offset(offset).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *gormAccountRepository) CountAccounts() (int64, error) {
	var count int64
	if err := r.db.Model(&entities.AccountEntity{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *gormAccountRepository) GetAccountByID(accountID int) (*entities.AccountEntity, error) {
	var account entities.AccountEntity
	if err := r.db.Preload("User").First(&account, accountID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrNotFound
		}
		return nil, err
	}
	return &account, nil
}

func (r *gormAccountRepository) GetAccountByPublicAddress(address string) (*entities.AccountEntity, error) {
	var account entities.AccountEntity
	if err := r.db.Where("public_address = ?", address).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrNotFound
		}
		return nil, err
	}
	return &account, nil
}

func (r *gormAccountRepository) UpdateAccount(txCtx context.Context, id int, data *entities.AccountEntity) error {
	tx := db.GetTx(txCtx, r.db)
	if err := tx.Model(&entities.AccountEntity{}).Where("id = ?", id).Select("*").Updates(data).Error; err != nil {
		return err
	}
	return nil
}
