package service

import "github.com/Narutchai01/solpay-core-service/internal/entities"

type ExampleService interface {
	GetExampleByID(id int) (entities.ExampleEntity, error)
}
