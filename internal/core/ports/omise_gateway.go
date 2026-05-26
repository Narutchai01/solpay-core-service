package ports

import (
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/omise/omise-go"
)

// OmiseGateway defines the interface for interacting with the Omise payment gateway.
type OmiseGateway interface {
	CreateRecipient(recipient request.CreateRecipient) (*omise.Recipient, error)
	CreateTransfer(amountSatang int64, recipientID string) (*omise.Transfer, error)
}
