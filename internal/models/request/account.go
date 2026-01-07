package request

type CreateAccountRequest struct {
	PublicAddress string `json:"public_address" binding:"required"`
	KycToken      string `json:"kyc_token"`
}
