package payment

import (
	"fmt"
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

// NewOmiseGateway creates a new OmiseGateway backed by an Omise client.
func NewOmiseGateway(client *omise.Client) ports.OmiseGateway {
	return &omiseGateway{client: client}
}

func (g *omiseGateway) CreateRecipient(recipientData request.CreateRecipient) (*omise.Recipient, error) {
	cfg := config.LoadConfig()
	recipient := &omise.Recipient{}

	op := &operations.CreateRecipient{
		Name: recipientData.Name,
		Type: omise.Individual,
		BankAccount: &omise.BankAccountRequest{
			Brand:  string(BankKBank),
			Number: recipientData.Number,
			Name:   recipientData.Name,
		},
	}

	if err := g.client.Do(recipient, op); err != nil {
		return nil, fmt.Errorf("create omise recipient: %w", err)
	}

	if err := g.verifyRecipient(cfg.OMISE_SECRET, recipient.ID); err != nil {
		return nil, fmt.Errorf("verify omise recipient: %w", err)
	}

	return recipient, nil
}

func (g *omiseGateway) CreateTransfer(amountSatang int64, recipientID string) (*omise.Transfer, error) {
	cfg := config.LoadConfig()
	transfer := &omise.Transfer{}

	op := &operations.CreateTransfer{
		Amount:    amountSatang + 3000,
		Recipient: recipientID,
		FailFast:  true,
	}

	if err := g.client.Do(transfer, op); err != nil {
		return nil, fmt.Errorf("create omise transfer: %w", err)
	}

	for _, action := range []string{"mark_as_sent", "mark_as_paid"} {
		if err := g.postTransferAction(cfg.OMISE_SECRET, transfer.ID, action); err != nil {
			return nil, fmt.Errorf("omise transfer %s: %w", action, err)
		}
	}

	return transfer, nil
}

// verifyRecipient sends a PATCH request to verify an Omise recipient.
func (g *omiseGateway) verifyRecipient(secret, recipientID string) error {
	req, err := http.NewRequest("PATCH",
		"https://api.omise.co/recipients/"+recipientID+"/verify",
		nil)
	if err != nil {
		return fmt.Errorf("build verify request: %w", err)
	}
	req.SetBasicAuth(secret, "")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute verify request: %w", err)
	}
	resp.Body.Close()
	return nil
}

// postTransferAction sends a POST request for a transfer action (e.g. mark_as_sent, mark_as_paid).
func (g *omiseGateway) postTransferAction(secret, transferID, action string) error {
	req, err := http.NewRequest("POST",
		"https://api.omise.co/transfers/"+transferID+"/"+action,
		nil)
	if err != nil {
		return fmt.Errorf("build %s request: %w", action, err)
	}
	req.SetBasicAuth(secret, "")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute %s request: %w", action, err)
	}
	resp.Body.Close()
	return nil
}
