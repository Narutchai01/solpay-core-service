package repositories

import (
	"context"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/db"
	"gorm.io/gorm"
)

type sqlUnitOfWork struct {
	db *gorm.DB
}

// NewSqlUnitOfWork creates a new UnitOfWork backed by GORM transactions.
func NewSqlUnitOfWork(database *gorm.DB) ports.UnitOfWork {
	return &sqlUnitOfWork{db: database}
}

func (u *sqlUnitOfWork) Execute(ctx context.Context, fn func(ctx context.Context) (any, error)) (any, error) {
	var result any
	err := u.db.Transaction(func(tx *gorm.DB) error {
		txCtx := db.NewTxContext(ctx, tx)

		res, err := fn(txCtx)
		if err != nil {
			return err
		}

		result = res
		return nil
	})

	return result, err
}
