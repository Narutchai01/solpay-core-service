package services

import (
	"errors"
	"fmt"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
)

// ExampleService defines operations for the example domain.
type ExampleService interface {
	GetExampleByID(id int) (entities.ExampleEntity, error)
}

type exampleService struct {
	repo ports.ExampleRepository
}

// NewExampleService creates a new ExampleService.
func NewExampleService(r ports.ExampleRepository) ExampleService {
	return &exampleService{repo: r}
}

func (s *exampleService) GetExampleByID(id int) (entities.ExampleEntity, error) {
	example, err := s.repo.GetExampleByID(id)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			return entities.ExampleEntity{}, entities.NewAppError(entities.ErrTypeNotFound, fmt.Sprintf("example %d not found", id), err)
		}
		return entities.ExampleEntity{}, entities.NewAppError(entities.ErrTypeInternal, "failed to get example", err)
	}
	return example, nil
}
