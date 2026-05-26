package response

type QuoteResponse struct {
	QuoteID          string   `json:"quote_id"`
	THBAmount        Decimal2 `json:"thb_amount"`
	USDTAmount       float64  `json:"usdt_amount"`
	ExchangeRate     float64  `json:"exchange_rate"`
	PromptPayID      *string  `json:"promptpay_id,omitempty"`
	Fee              float64  `json:"fee"`
	QuoteTpye        string   `json:"quote_type"`
	ExpiresInSeconds int      `json:"expires_in_seconds"`
}
