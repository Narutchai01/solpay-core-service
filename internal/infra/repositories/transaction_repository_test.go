package repositories_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Narutchai01/solpay-core-service/internal/infra/repositories"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGormTransactionRepositoryGetMonthlySpendingSummary(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	ctx := context.Background()

	repo := repositories.NewGormTransactionRepository(gormDB)

	t.Run("GetMonthlySpendingSummary_Success", func(t *testing.T) {
		accountID := uint(1)
		limit := 6

		mock.ExpectQuery(`SELECT TO_CHAR\(created_at, 'YYYY-MM'\) AS month, COALESCE\(SUM\(thb_amount\), 0\) AS total_spent FROM transaction_entities WHERE account_id = \$1 AND deleted_at IS NULL AND status = 'COMPLETED' AND transaction_type IN \('OFFCHAIN', 'ONCHAIN'\) GROUP BY month ORDER BY month DESC LIMIT \$2`).
			WithArgs(accountID, limit).
			WillReturnRows(sqlmock.NewRows([]string{"month", "total_spent"}).
				AddRow("2026-05", 120000.0).
				AddRow("2026-04", 98000.0))

		summaries, err := repo.GetMonthlySpendingSummary(ctx, accountID, limit)
		assert.NoError(t, err)
		assert.Len(t, summaries, 2)
		assert.Equal(t, "2026-05", summaries[0].Month)
		assert.Equal(t, 120000.0, summaries[0].TotalSpent)
		assert.Equal(t, "2026-04", summaries[1].Month)
		assert.Equal(t, 98000.0, summaries[1].TotalSpent)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
