package repositories

import "github.com/Narutchai01/solpay-core-service/internal/entities"

type ExampleRepository interface {
	GetExampleByID(id int) (entities.ExampleEntity, error)
}
