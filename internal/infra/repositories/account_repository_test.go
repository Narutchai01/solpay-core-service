package repositories_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/infra/repositories"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGormAccountRepositoryCreateAccount(t *testing.T) {
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

	repo := repositories.NewGormAccountRepository(gormDB)

	t.Run("CreateAccount_Success", func(t *testing.T) {
		expectData := entities.AccountEntity{
			PublicAddress: "test_public_address",
			KycToken:      nil,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "account_entities"`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), expectData.PublicAddress).
			WillReturnRows(sqlmock.NewRows([]string{"id", "kyc_token"}).AddRow(1, nil))
		mock.ExpectCommit()

		err := repo.CreateAccount(ctx, &expectData)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CreateAccount_DuplicateKeyError", func(t *testing.T) {
		expectData := entities.AccountEntity{
			PublicAddress: "test_public_address",
			KycToken:      nil,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "account_entities"`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), expectData.PublicAddress).
			WillReturnError(gorm.ErrDuplicatedKey)
		mock.ExpectRollback()

		err := repo.CreateAccount(ctx, &expectData)
		assert.Equal(t, entities.ErrConflict, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CreateAccount_OtherError", func(t *testing.T) {
		expectData := entities.AccountEntity{
			PublicAddress: "test_public_address",
			KycToken:      nil,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "account_entities"`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), expectData.PublicAddress).
			WillReturnError(gorm.ErrInvalidData)
		mock.ExpectRollback()

		err := repo.CreateAccount(ctx, &expectData)
		assert.Equal(t, gorm.ErrInvalidData, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
func TestGormAccountRepositoryGetAccounts(t *testing.T) {
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
	t.Run("GetAccounts_Success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "account_entities" WHERE "account_entities"\."deleted_at" IS NULL LIMIT \$1`).
			WithArgs(10).
			WillReturnRows(sqlmock.NewRows([]string{"id", "public_address", "kyc_token"}).
				AddRow(1, "address1", nil).
				AddRow(2, "address2", nil))

		accounts, err := repo.GetAccounts(1, 10)
		assert.NoError(t, err)
		assert.Len(t, accounts, 2)
		assert.Equal(t, "address1", accounts[0].PublicAddress)
		assert.Equal(t, "address2", accounts[1].PublicAddress)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetAccounts_EmptyResult", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "account_entities" WHERE "account_entities"\."deleted_at" IS NULL LIMIT \$1 OFFSET \$2`).
			WithArgs(10, 10).
			WillReturnRows(sqlmock.NewRows([]string{"id", "public_address", "kyc_token"}))

		accounts, err := repo.GetAccounts(2, 10)
		assert.NoError(t, err)
		assert.Len(t, accounts, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetAccounts_DatabaseError", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "account_entities" WHERE "account_entities"\."deleted_at" IS NULL LIMIT \$1`).
			WithArgs(10).
			WillReturnError(gorm.ErrInvalidData)

		accounts, err := repo.GetAccounts(1, 10)
		assert.Error(t, err)
		assert.Nil(t, accounts)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGormAccountRepositoryCountAccounts(t *testing.T) {
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

	t.Run("CountAccounts_Success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT count\(\*\) FROM "account_entities" WHERE "account_entities"\."deleted_at" IS NULL`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

		count, err := repo.CountAccounts()
		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CountAccounts_DatabaseError", func(t *testing.T) {
		mock.ExpectQuery(`SELECT count\(\*\) FROM "account_entities" WHERE "account_entities"\."deleted_at" IS NULL`).
			WillReturnError(gorm.ErrInvalidData)

		count, err := repo.CountAccounts()
		assert.Error(t, err)
		assert.Equal(t, int64(0), count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGormAccountRepositoryGetAccountByID(t *testing.T) {
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

	t.Run("GetAccountByID_Success", func(t *testing.T) {
		now := time.Now()
		mock.ExpectQuery(`SELECT \* FROM "account_entities" WHERE "account_entities"\."id" = \$1 AND "account_entities"\."deleted_at" IS NULL ORDER BY "account_entities"\."id" LIMIT \$2`).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "public_address", "kyc_token", "created_at", "updated_at", "deleted_at"}).
				AddRow(1, "address1", nil, now, now, nil))

		account, err := repo.GetAccountByID(1)
		assert.NoError(t, err)
		assert.Equal(t, "address1", account.PublicAddress)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetAccountByID_NotFound", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "account_entities" WHERE "account_entities"\."id" = \$1 AND "account_entities"\."deleted_at" IS NULL ORDER BY "account_entities"\."id" LIMIT \$2`).
			WithArgs(2, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		account, err := repo.GetAccountByID(2)
		assert.Equal(t, entities.ErrNotFound, err)
		assert.Equal(t, &entities.AccountEntity{}, account)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetAccountByID_DatabaseError", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "account_entities" WHERE "account_entities"\."id" = \$1 AND "account_entities"\."deleted_at" IS NULL ORDER BY "account_entities"\."id" LIMIT \$2`).
			WithArgs(3, 1).
			WillReturnError(gorm.ErrInvalidData)

		account, err := repo.GetAccountByID(3)
		assert.Equal(t, entities.ErrInternal, err)
		assert.Equal(t, &entities.AccountEntity{}, account)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
