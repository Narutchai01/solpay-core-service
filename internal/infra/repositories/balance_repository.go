package repositories

import (
	"context"
	"errors"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/db"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"gorm.io/gorm"
)

type balanceRepository struct {
	db *gorm.DB
}

// NewGormBalanceRepository creates a new GORM-backed BalanceRepository.
func NewGormBalanceRepository(database *gorm.DB) ports.BalanceRepository {
	return &balanceRepository{db: database}
}

func (r *balanceRepository) CreateBalance(txCtx context.Context, data *entities.BalanceEntity) error {
	tx := db.GetTx(txCtx, r.db)
	return tx.Create(data).Error
}

func (r *balanceRepository) GetBalances(page, limit int) ([]entities.BalanceEntity, error) {
	var balances []entities.BalanceEntity
	offset := (page - 1) * limit

	if err := r.db.Limit(limit).Offset(offset).Find(&balances).Error; err != nil {
		return nil, err
	}
	return balances, nil
}

func (r *balanceRepository) CountBalances() (int64, error) {
	var count int64
	if err := r.db.Model(&entities.BalanceEntity{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *balanceRepository) GetBalanceByID(balanceID int) (*entities.BalanceEntity, error) {
	var balance entities.BalanceEntity
	if err := r.db.First(&balance, balanceID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrNotFound
		}
		return nil, err
	}
	return &balance, nil
}

func (r *balanceRepository) UpdateBalance(txCtx context.Context, data *entities.BalanceEntity) error {
	tx := db.GetTx(txCtx, r.db)
	return tx.Model(&entities.BalanceEntity{}).
		Where("account_id = ?", data.AccountID).
		Select("*").
		Updates(data).Error
}

func (r *balanceRepository) GetBalanceByAccountID(accountID uint) (*entities.BalanceEntity, error) {
	var balance entities.BalanceEntity
	if err := r.db.Where("account_id = ?", accountID).First(&balance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrNotFound
		}
		return nil, err
	}
	return &balance, nil
}
