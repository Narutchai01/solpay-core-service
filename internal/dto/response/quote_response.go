package response

type QuoteResponse struct {
	QuoteID          string  `json:"quote_id"`
	THBAmount        float64 `json:"thb_amount"`
	USDTAmount       float64 `json:"usdt_amount"`
	ExchangeRate     float64 `json:"exchange_rate"`
	Fee              float64 `json:"fee"`
	ExpiresInSeconds int     `json:"expires_in_seconds"`
}
