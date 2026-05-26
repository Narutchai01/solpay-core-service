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

type TransactionQuery struct {
	Page     int      `query:"page"`
	PageSize int      `query:"pageSize"`
	TxType   []string `query:"txType"`
}

func (q *TransactionQuery) GetOffset() int {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
	return (q.Page - 1) * q.PageSize
}

func (q *TransactionQuery) GetLimit() int {
	if q.PageSize <= 0 {
		return 10
	}
	return q.PageSize
}

type QueryTransactionSummaryRequest struct {
	Month string `json:"month"`
	Year  string `json:"year"`
}
