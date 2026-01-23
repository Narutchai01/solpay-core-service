package repositories_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/infra/repositories"
	assert "github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGormAccountRepository_CreateAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	repo := repositories.NewGormAccountRepository(gormDB)

	expectData := entities.AccountEntity{
		PublicAddress: "test_public_address",
		KycToken:      nil,
	}

	t.Run("CreateAccount_Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "account_entities"`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), expectData.PublicAddress, expectData.KycToken).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		acc, err := repo.CreateAccount(nil, &expectData)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, &expectData.PublicAddress, &acc.PublicAddress)
		assert.NotNil(t, acc.ID, "Account ID should not be nil")
	})
}
