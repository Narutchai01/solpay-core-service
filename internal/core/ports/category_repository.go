package ports

import "github.com/Narutchai01/solpay-core-service/internal/entities"

type CategoryRepository interface {
	GetCategories() ([]entities.Category, error)
	GetCategory(id int) (*entities.Category, error)
}
