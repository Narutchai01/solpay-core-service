package repositories

import (
	"errors"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"gorm.io/gorm"
)

type gormExampleRepository struct {
	db *gorm.DB
}

// NewGormExampleRepository creates a new GORM-backed ExampleRepository.
func NewGormExampleRepository(database *gorm.DB) ports.ExampleRepository {
	return &gormExampleRepository{db: database}
}

func (r *gormExampleRepository) GetExampleByID(id int) (entities.ExampleEntity, error) {
	var example entities.ExampleEntity
	if err := r.db.First(&example, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.ExampleEntity{}, entities.ErrNotFound
		}
		return entities.ExampleEntity{}, err
	}
	return example, nil
}
