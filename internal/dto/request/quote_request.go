package request

type CreateQuoteRequest struct {
	THBAmount  float64 `json:"thb_amount"`  // รับมาเป็นบาท เช่น 1000.00
	ActionType string  `json:"action_type"` // เช่น "TOPUP_CRYPTO"
}
