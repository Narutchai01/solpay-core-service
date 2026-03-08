package ports

type Consumer interface {
	TransactionOrchestrator() error
	BalanceConsumer() error
	PaymentConsumer() error
}
