package request

type CreateTransactionRequest struct {
	TransactionType string `json:"transaction_type" validate:"required,oneof=top_up transaction_onchain transaction_offchain"`

	THBAmount  float64 `json:"thb_amount" validate:"min=0"`
	USDTAmount float64 `json:"usdt_amount" validate:"min=0"`

	Fee float64 `json:"fee" validate:"min=0"`

	PromptPayID *string `json:"propmtpayID"`

	TxHash *string `json:"tx_hash"`
}
