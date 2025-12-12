package services

import (
	ports "github.com/Narutchai01/solpay-core-service/internal/core/ports/repositories"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
)

type ExampleService struct {
	repo ports.ExampleRepository
}

func NewExampleService(r ports.ExampleRepository) *ExampleService {
	return &ExampleService{repo: r}
}

func (s *ExampleService) GetExampleByID(id int) (entities.ExampleEntity, error) {
	return s.repo.GetExampleByID(id)
}
