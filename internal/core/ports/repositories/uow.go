package repositories

import "context"

// domain/repository.go
type UnitOfWork interface {
	// Execute รับ function ที่เราต้องการให้รันใน Transaction
	Execute(ctx context.Context, fn func(ctx context.Context) (any, error)) (any, error)
}
