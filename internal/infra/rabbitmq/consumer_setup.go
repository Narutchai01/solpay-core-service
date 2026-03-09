package rabbitmq

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/infra/payment"
	"github.com/Narutchai01/solpay-core-service/internal/infra/repositories"
	"github.com/Narutchai01/solpay-core-service/internal/websocket"
	"github.com/omise/omise-go"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type ConsumerSetup struct {
	channel *amqp.Channel
	db      *gorm.DB
	wsHub   *websocket.Hub
	omise   *omise.Client
}

func NewConsumerSetup(channel *amqp.Channel, db *gorm.DB, wsHub *websocket.Hub, omise *omise.Client) *ConsumerSetup {
	return &ConsumerSetup{
		channel: channel,
		db:      db,
		wsHub:   wsHub,
		omise:   omise,
	}
}

func (cs *ConsumerSetup) Setup() {

	transactionRepo := repositories.NewGormTransactionRepository(cs.db)
	balanceRepo := repositories.NewGormBalanceRepository(cs.db)
	paymentRepo := repositories.NewGormPaymentRepository(cs.db)
	uowRepo := repositories.NewSqlUnitOfWork(cs.db)
	Publisher := NewPublisher(cs.channel)
	omiseGateway := payment.NewOmiseGateway(cs.omise)
	transactionService := services.NewTransactionService(transactionRepo, uowRepo, Publisher, cs.wsHub)
	balanceService := services.NewBalanceService(balanceRepo, uowRepo, Publisher)
	paymentService := services.NewPaymentService(omiseGateway, paymentRepo)
	consumer := NewConsumer(cs.channel, transactionService, balanceService, paymentService, Publisher)

	go func() {
		if err := consumer.TransactionOrchestrator(); err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := consumer.BalanceConsumer(); err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := consumer.PaymentConsumer(); err != nil {
			panic(err)
		}
	}()
}
