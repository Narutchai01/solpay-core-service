package request

type CreateAccountRequest struct {
	PublicAddress string `json:"public_address" validate:"required"`
	KycToken      string `json:"kyc_token"`
}

type GetAccountsRequest struct {
	Page  int `json:"page" validate:"required,min=1"`
	Limit int `json:"limit" validate:"required,min=1,max=100"`
}
