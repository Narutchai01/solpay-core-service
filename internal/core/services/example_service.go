package services

import (
	"errors"

	ports "github.com/Narutchai01/solpay-core-service/internal/core/ports/repositories"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
)

type ExampleService interface {
	GetExampleByID(id int) (entities.ExampleEntity, error)
}

type exampleService struct {
	repo ports.ExampleRepository
}

func NewExampleService(r ports.ExampleRepository) ExampleService {
	return &exampleService{repo: r}
}

func (s *exampleService) GetExampleByID(id int) (entities.ExampleEntity, error) {
	var example entities.ExampleEntity
	example, err := s.repo.GetExampleByID(id)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			return entities.ExampleEntity{}, entities.NewAppError(entities.ErrTypeNotFound, "example not found", err)
		}
		return entities.ExampleEntity{}, entities.NewAppError(entities.ErrTypeInternal, "internal server error", err)
	}
	return example, nil
}
