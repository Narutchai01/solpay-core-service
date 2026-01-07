package services

import (
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
	return s.repo.GetExampleByID(id)
}
