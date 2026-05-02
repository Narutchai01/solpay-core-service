package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost                         string
	DBPort                         string
	DBUser                         string
	DBPassword                     string
	DBName                         string
	APPPort                        string
	TimeZone                       string
	Environment                    string
	SECRET_JWT                     string
	JWT_EXPIRATION_HOURS           string
	RABBITMQ_URL                   string
	SOLANA_WORK_QUEUE              string
	TRANSACTION_ORCHESTRATOR_QUEUE string
	BALANCE_QUEUE                  string
	PAYMENT_QUEUE                  string
	OMISE_KEY                      string
	OMISE_SECRET                   string
	MINT_TOKEN_ADDRESS             string
	RECEIVE_ADDRESS                string
	RPC_URL                        string
	SUPABASE_PRIVATE_KEY           string
	SUPABASE_URL                   string
	SWAP_SERVICE_URL               string
}

func GetEnv(key string, fallback ...string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	if len(fallback) > 0 {
		return fallback[0]
	}
	msg := fmt.Sprintf("Environment variable %s is required but not set", key)
	panic(msg)
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: .env file not found. Using system environment variables instead.")
	}
	return &Config{
		DBHost:                         GetEnv("DB_HOST"),
		DBPort:                         GetEnv("DB_PORT"),
		DBUser:                         GetEnv("DB_USER"),
		DBPassword:                     GetEnv("DB_PASSWORD"),
		DBName:                         GetEnv("DB_NAME"),
		APPPort:                        GetEnv("APP_PORT", "8080"),
		TimeZone:                       GetEnv("TIMEZONE", "Asia/Bangkok"),
		Environment:                    GetEnv("ENVIRONMENT", "development"),
		SECRET_JWT:                     GetEnv("SECRET_JWT"),
		JWT_EXPIRATION_HOURS:           GetEnv("JWT_EXPIRATION_HOURS", "72"),
		RABBITMQ_URL:                   GetEnv("MQ_URL", "amqp://guest:guest@localhost:5672/"),
		SOLANA_WORK_QUEUE:              GetEnv("SOLANA_WORK_QUEUE", "solana-worker.tx.submit"),
		TRANSACTION_ORCHESTRATOR_QUEUE: GetEnv("CORE_TRANSACTION_QUEUE", "core.transaction.status.update"),
		BALANCE_QUEUE:                  GetEnv("CORE_BALANCE_QUEUE", "core.balance.update"),
		PAYMENT_QUEUE:                  GetEnv("CORE_PAYMENT_QUEUE", "core.payment.update"),
		OMISE_KEY:                      GetEnv("OMISE_KEY"),
		OMISE_SECRET:                   GetEnv("OMISE_SECRET"),
		MINT_TOKEN_ADDRESS:             GetEnv("MINT_TOKEN_ADDRESS"),
		RECEIVE_ADDRESS:                GetEnv("RECEIVE_ADDRESS"),
		RPC_URL:                        GetEnv("RPC_URL", "https://api.devnet.solana.com"),
		SUPABASE_PRIVATE_KEY:           GetEnv("SUPABASE_PRIVATE_KEY"),
		SUPABASE_URL:                   GetEnv("SUPABASE_URL"),
		SWAP_SERVICE_URL:               GetEnv("SWAP_SERVICE_URL", "http://localhost:8081"),
	}
}
