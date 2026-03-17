package request

type OffChainRequest struct {
	PromptPayID string `json:"promptpay_id" validate:"required"`
	THBAmount   int64  `json:"thb_amount" validate:"required,gt=0"`
}
