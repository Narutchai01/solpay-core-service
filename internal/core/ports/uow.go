package ports

import "context"

// UnitOfWork executes a function within a database transaction.
type UnitOfWork interface {
	Execute(ctx context.Context, fn func(ctx context.Context) (any, error)) (any, error)
}
