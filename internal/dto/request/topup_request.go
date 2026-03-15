package request

type TopUpRequest struct {
	QuoteID string `json:"quote_id" validate:"required"`
	TxHash  string `json:"tx_hash" validate:"required"`
}
