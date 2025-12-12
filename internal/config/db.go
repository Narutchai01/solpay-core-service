package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		panic("DB_HOST environment variable is required")
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		panic("DB_PORT environment variable is required")
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		panic("DB_USER environment variable is required")
	}
	pass := os.Getenv("DB_PASS")
	if pass == "" {
		panic("DB_PASS environment variable is required")
	}
	name := os.Getenv("DB_NAME")
	if name == "" {
		panic("DB_NAME environment variable is required")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Bangkok", host, port, user, pass, name)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	db.AutoMigrate(entities.ExampleEntity{})

	if err != nil {
		return nil, err
	}

	return db, nil

}
