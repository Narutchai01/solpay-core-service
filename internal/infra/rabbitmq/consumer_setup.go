package rabbitmq

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/infra/repositories"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type ConsumerSetup struct {
	channel *amqp.Channel
	db      *gorm.DB
}

func NewConsumerSetup(channel *amqp.Channel, db *gorm.DB) *ConsumerSetup {
	return &ConsumerSetup{
		channel: channel,
		db:      db,
	}
}

func (cs *ConsumerSetup) Setup() {

	go func() {
		transactionRepo := repositories.NewGormTransactionRepository(cs.db)
		balanceRepo := repositories.NewGormBalanceRepository(cs.db)
		uowRepo := repositories.NewSqlUnitOfWork(cs.db)
		Publisher := NewPublisher(cs.channel)
		transactionService := services.NewTransactionService(transactionRepo, uowRepo, Publisher)
		balanceService := services.NewBalanceService(balanceRepo, uowRepo)
		consumer := NewConsumer(cs.channel, transactionService, balanceService)
		if err := consumer.TransactionOrchestrator(); err != nil {
			panic(err)
		}
	}()
}
