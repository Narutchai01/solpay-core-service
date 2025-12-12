package repositories

import (
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
		return entities.ExampleEntity{}, err
	}
	return example, nil
}
