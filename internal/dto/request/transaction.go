package request

type CreateTransactionRequest struct {
	TransactionType string

	THBAmount  float64
	USDTAmount float64

	Fee float64

	PromptPayID *string

	TxHash *string
}

type TransactionMessage struct {
	TxID         string `json:"tx_id"`
	SourceWorker string `json:"source_worker"`
	Status       string `json:"status"`
}

type QueryTransactionSummaryRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
