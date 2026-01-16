package request

type CreateTransactionRequest struct {
	AccountID  uint   `json:"account_id" validate:"required"`
	CategoryID uint   `json:"category_id" validate:"required"`
	Type       string `json:"type" validate:"required"`
	Status     string `json:"status" validate:"required"`
	THBAmount  int64  `json:"thb_amount" validate:"required,min=0"`
	USDTAmount int64  `json:"usdt_amount" validate:"required,min=0"`
	Fee        int64  `json:"fee"`
}
