package payment

import (
	"net/http"

	"github.com/Narutchai01/solpay-core-service/internal/config"
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
	cfg := config.LoadConfig()
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

	req, _ := http.NewRequest("PATCH",
		"https://api.omise.co/recipients/"+recipient.ID+"/verify",
		nil)
	req.SetBasicAuth(cfg.OMISE_SECRET, "")
	_, err = http.DefaultClient.Do(req)

	return recipient, nil
}

func (g *omiseGateway) CreateTransfer(amountSatang int64, recipientID string) (*omise.Transfer, error) {
	cfg := config.LoadConfig()
	transfer := &omise.Transfer{}

	op := &operations.CreateTransfer{
		Amount:    (amountSatang * 100) + 3000,
		Recipient: recipientID,
		FailFast:  true,
	}

	err := g.client.Do(transfer, op)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest("POST",
		"https://api.omise.co/transfers/"+transfer.ID+"/mark_as_sent",
		nil)
	req.SetBasicAuth(cfg.OMISE_SECRET, "")
	http.DefaultClient.Do(req)

	req, _ = http.NewRequest("POST",
		"https://api.omise.co/transfers/"+transfer.ID+"/mark_as_paid",
		nil)
	req.SetBasicAuth(cfg.OMISE_SECRET, "")
	http.DefaultClient.Do(req)

	return transfer, nil
}
