package request

type GetBalancesRequest struct {
	Page  int `json:"page" validate:"required,min=1"`
	Limit int `json:"limit" validate:"required,min=1,max=100"`
}

type UpdateBalanceCommand struct {
	TransactionID string `json:"transaction_id"`
	AccountID     uint   `json:"account_id"`
	Action        string `json:"action"`
	THBAmount     int64  `json:"thb_amount"`
	USDTAmount    int64  `json:"usdt_amount"`
}
