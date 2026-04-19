package services

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
)

type CategoryService interface {
	GetCategories() ([]entities.Category, error)
}

type categoryService struct {
	categoryRepo ports.CategoryRepository
}

func NewCategoryService(categoryRepo ports.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func (s *categoryService) GetCategories() ([]entities.Category, error) {
	return s.categoryRepo.GetCategories()
}
