package repositories

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"gorm.io/gorm"
)

type gormQuoteRepository struct {
	db *gorm.DB
}

func NewGormQuoteRepository(database *gorm.DB) ports.QuoteRepository {
	return &gormQuoteRepository{db: database}
}

func (r *gormQuoteRepository) CreateQuote(quote *entities.Quote) error {
	if err := r.db.Create(quote).Error; err != nil {
		return err
	}
	return nil
}

func (r *gormQuoteRepository) GetQuoteByID(quoteID string) (*entities.Quote, error) {
	var quote entities.Quote
	if err := r.db.Where("id = ?", quoteID).First(&quote).Error; err != nil {
		return nil, err
	}
	return &quote, nil
}

func (r *gormQuoteRepository) UpdateQuote(quote *entities.Quote) error {
	if err := r.db.Save(quote).Error; err != nil {
		return err
	}
	return nil
}
