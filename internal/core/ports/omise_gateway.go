package ports

import (
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/omise/omise-go"
)

type OmiseGateway interface {
	CreateRecipient(recipent request.CreateRecipient) (*omise.Recipient, error)
	CreateTransfer(amountSatang int64, recipientID string) (*omise.Transfer, error)
}
