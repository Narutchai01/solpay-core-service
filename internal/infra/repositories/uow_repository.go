package repositories

import (
	"context"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports/repositories"
	"gorm.io/gorm"
)

type SqlUnitOfWork struct {
	db *gorm.DB
}

func NewSqlUnitOfWork(db *gorm.DB) repositories.UnitOfWork {
	return &SqlUnitOfWork{db: db}
}

func (u *SqlUnitOfWork) Execute(ctx context.Context, fn func(ctx context.Context) (any, error)) (any, error) {
	var result any
	err := u.db.Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, "tx_key", tx)

		res, err := fn(txCtx) // รับค่าที่ return จาก UseCase
		if err != nil {
			return err
		}

		result = res // เก็บค่าไว้ในตัวแปรด้านนอก
		return nil
	})

	return result, err
}
