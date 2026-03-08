package payment

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

type omiseGateway struct {
	client *omise.Client
}

func NewOmiseGateway(client *omise.Client) ports.OmiseGateway {
	return &omiseGateway{client: client}
}

func (g *omiseGateway) CreateRecipient(recinpientData request.CreateRecipient) (*omise.Recipient, error) {
	recipient := &omise.Recipient{}

	op := &operations.CreateRecipient{
		Name: recinpientData.Name, // ในอนาคตอาจจะรับมาจาก input
		Type: omise.Individual,    // ล็อกเป็นบุคคลธรรมดาไว้ก่อน
		BankAccount: &omise.BankAccountRequest{
			Brand:  string(BankKBank),
			Number: recinpientData.Number,
			Name:   recinpientData.Name,
		},
	}

	err := g.client.Do(recipient, op)
	if err != nil {
		return nil, err
	}

	return recipient, nil
}

func (g *omiseGateway) CreateTransfer(amountSatang int64, recipientID string) (*omise.Transfer, error) {
	transfer := &omise.Transfer{}

	op := &operations.CreateTransfer{
		Amount:    amountSatang,
		Recipient: recipientID,
		FailFast:  true,
	}

	err := g.client.Do(transfer, op)
	if err != nil {
		return nil, err
	}

	return transfer, nil
}
