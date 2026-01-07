package repositories

import (
	"errors"

	ports "github.com/Narutchai01/solpay-core-service/internal/core/ports/repositories"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"gorm.io/gorm"
)

type GormExampleRepository struct {
	db *gorm.DB
}

func NewGormExampleRepository(db *gorm.DB) ports.ExampleRepository {
	return &GormExampleRepository{db: db}
}

func (r *GormExampleRepository) GetExampleByID(id int) (entities.ExampleEntity, error) {
	var example entities.ExampleEntity
	if err := r.db.First(&example, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.ExampleEntity{}, entities.ErrNotFound
		}
		return entities.ExampleEntity{}, entities.ErrInternal
	}
	return example, nil
}
