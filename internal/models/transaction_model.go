package models

type CreateTransactionTopUp struct {
	TransactionType string
	THBAmount       float64
	USDTAmount      float64
	Fee             float64
	TxHash          *string
}
