package ports

type Consumer interface {
	TransactionOrchestrator() error
}
