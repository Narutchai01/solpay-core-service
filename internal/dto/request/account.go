package request

type CreateAccountRequest struct {
	PublicAddress string `json:"public_address" validate:"required"`
	KycToken      string `json:"kyc_token"`
}
