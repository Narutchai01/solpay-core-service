package db

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

// txKey is an unexported type used as the context key for database transactions.
type txKey struct{}

func ConnectDB() (*gorm.DB, error) {
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Bangkok",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:         newLogger,
		TranslateError: true,
	})
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	slog.Info("Database connected successfully")

	if err := db.AutoMigrate(
		&entities.ExampleEntity{},
		&entities.AccountEntity{},
		&entities.BalanceEntity{},
		&entities.TransactionEntity{},
		&entities.TransactionOnChain{},
		&entities.TransactionOffChain{},
		&entities.Recipient{},
		&entities.LogPayment{},
		&entities.Quote{},
		&entities.AdminEntity{},
		&entities.Category{},
	); err != nil {
		return nil, fmt.Errorf("auto-migrate: %w", err)
	}

	if err := seedCategories(db); err != nil {
		return nil, err
	}

	return db, nil
}

func seedCategories(db *gorm.DB) error {
	categories := []entities.Category{
		{ID: 1, Name: "Others"},
		{ID: 2, Name: "Food/Drink"},
		{ID: 3, Name: "Shopping"},
		{ID: 4, Name: "invest"},
	}

	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&categories).Error; err != nil {
		return fmt.Errorf("seed categories: %w", err)
	}

	return nil
}

// NewTxContext stores the transaction in the context.
func NewTxContext(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// GetTx retrieves the transaction from the context, falling back to defaultDB.
func GetTx(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return defaultDB
}
