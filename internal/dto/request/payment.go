package request

type CreateRecipient struct {
	Name   string `json:"name" validate:"required"`
	Number string `json:"number" validate:"required"`
}

type RequestPaymentQueue struct {
	TransactionID string `json:"transaction_id" validate:"required"`
	Amount        int64  `json:"amount" validate:"required"`
	Number        string `json:"number" validate:"required"`
}
