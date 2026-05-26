package rabbitmq

import (
	"log"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/infra/payment"
	"github.com/Narutchai01/solpay-core-service/internal/infra/repositories"
	"github.com/Narutchai01/solpay-core-service/internal/websocket"
	"github.com/omise/omise-go"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

// ConsumerSetup wires up all consumer dependencies and starts goroutines.
type ConsumerSetup struct {
	channel *amqp.Channel
	db      *gorm.DB
	wsHub   *websocket.Hub
	omise   *omise.Client
}

// NewConsumerSetup creates a new ConsumerSetup.
func NewConsumerSetup(channel *amqp.Channel, db *gorm.DB, wsHub *websocket.Hub, omise *omise.Client) *ConsumerSetup {
	return &ConsumerSetup{
		channel: channel,
		db:      db,
		wsHub:   wsHub,
		omise:   omise,
	}
}

// Setup initialises services, creates the consumer, and starts consumer goroutines.
func (cs *ConsumerSetup) Setup() {
	transactionRepo := repositories.NewGormTransactionRepository(cs.db)
	balanceRepo := repositories.NewGormBalanceRepository(cs.db)
	paymentRepo := repositories.NewGormPaymentRepository(cs.db)
	uow := repositories.NewSqlUnitOfWork(cs.db)
	pub := NewPublisher(cs.channel)
	omiseGW := payment.NewOmiseGateway(cs.omise)
	cfg := config.LoadConfig()

	transactionService := services.NewTransactionService(transactionRepo, uow, pub, cs.wsHub, cfg)
	balanceService := services.NewBalanceService(balanceRepo, uow, pub)
	paymentService := services.NewPaymentService(omiseGW, paymentRepo)

	c := NewConsumer(cs.channel, transactionService, balanceService, paymentService, pub)

	go func() {
		if err := c.TransactionOrchestrator(); err != nil {
			log.Fatalf("Transaction orchestrator consumer failed: %v", err)
		}
	}()

	go func() {
		if err := c.BalanceConsumer(); err != nil {
			log.Fatalf("Balance consumer failed: %v", err)
		}
	}()

	go func() {
		if err := c.PaymentConsumer(); err != nil {
			log.Fatalf("Payment consumer failed: %v", err)
		}
	}()
}
