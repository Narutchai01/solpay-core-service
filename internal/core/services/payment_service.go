package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/google/uuid"
)

// PaymentService defines the interface for payment processing operations.
type PaymentService interface {
	ProcessPayment(ctx context.Context, paymentQueue []byte) ([]byte, error)
}

type paymentService struct {
	gateway     ports.OmiseGateway
	paymentRepo ports.PaymentRepository
}

// NewPaymentService creates a new PaymentService.
func NewPaymentService(gateway ports.OmiseGateway, paymentRepo ports.PaymentRepository) PaymentService {
	return &paymentService{
		gateway:     gateway,
		paymentRepo: paymentRepo,
	}
}

func (s *paymentService) ProcessPayment(ctx context.Context, paymentQueue []byte) ([]byte, error) {
	var data request.RequestPaymentQueue
	if err := json.Unmarshal(paymentQueue, &data); err != nil {
		return nil, fmt.Errorf("unmarshal payment queue: %w", err)
	}

	recipient, err := s.findOrCreateRecipient(data)
	if err != nil {
		return nil, fmt.Errorf("find or create recipient: %w", err)
	}

	transfer, err := s.gateway.CreateTransfer(data.Amount, recipient.RecipientID)
	if err != nil {
		return nil, fmt.Errorf("create transfer: %w", err)
	}

	transactionUUID, err := uuid.Parse(data.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("parse transaction UUID %q: %w", data.TransactionID, err)
	}

	logPayment := entities.LogPayment{
		TransactionUUID:    transactionUUID,
		OmiseTransactionID: transfer.ID,
	}
	if err := s.paymentRepo.CreateLogPayment(&logPayment); err != nil {
		return nil, fmt.Errorf("create log payment: %w", err)
	}

	payload := request.TransactionMessage{
		TxID:         data.TransactionID,
		SourceWorker: "PAYMENT",
		Status:       string(entities.StatusPaymentSuccess),
	}

	return json.Marshal(payload)
}

// findOrCreateRecipient looks up an existing recipient or creates a new one.
func (s *paymentService) findOrCreateRecipient(data request.RequestPaymentQueue) (entities.Recipient, error) {
	recipient, err := s.paymentRepo.GetRecipientByNumber(data.Number)
	if err == nil {
		return recipient, nil
	}

	if !errors.Is(err, entities.ErrNotFound) {
		return entities.Recipient{}, fmt.Errorf("get recipient by number: %w", err)
	}

	// Recipient not found — create a new one via the payment gateway.
	rpReq := request.CreateRecipient{
		Name:   "John Doe",
		Number: data.Number,
	}

	rp, err := s.gateway.CreateRecipient(rpReq)
	if err != nil {
		return entities.Recipient{}, fmt.Errorf("create recipient via gateway: %w", err)
	}

	recipient = entities.Recipient{Number: data.Number, RecipientID: rp.ID}
	if err := s.paymentRepo.CreateRecipient(&recipient); err != nil {
		return entities.Recipient{}, fmt.Errorf("persist recipient: %w", err)
	}

	return recipient, nil
}
