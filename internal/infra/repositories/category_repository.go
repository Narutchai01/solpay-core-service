package repositories

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) ports.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) GetCategories() ([]entities.Category, error) {
	var categories []entities.Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
