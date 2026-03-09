package services

import (
	"context"
	"encoding/json"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/google/uuid"
)

type PaymentService interface {
	ProcessPayment(ctx context.Context, paymentQueue []byte) error
}

type paymentService struct {
	paymentGateWay ports.OmiseGateway
	paymentRepo    ports.PaymentRepository
}

func NewPaymentService(paymentGateWay ports.OmiseGateway, paymentRepo ports.PaymentRepository) PaymentService {
	return &paymentService{
		paymentGateWay: paymentGateWay,
		paymentRepo:    paymentRepo,
	}
}

func (s *paymentService) ProcessPayment(ctx context.Context, paymentQueue []byte) error {
	var paymentData request.RequestPaymentQueue
	if err := json.Unmarshal(paymentQueue, &paymentData); err != nil {
		return err
	}

	recipient, err := s.paymentRepo.GetRecipentByNumber(paymentData.Number)
	if err != nil {
		if err == entities.ErrNotFound {

			rpRaw := request.CreateRecipient{
				Name:   "Joh doe",
				Number: paymentData.Number,
			}

			rp, err := s.paymentGateWay.CreateRecipient(rpRaw)
			if err != nil {
				return err
			}

			recipient = entities.Recipient{Number: paymentData.Number, RecipientID: rp.ID}
			if err := s.paymentRepo.CreateRecipient(&recipient); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	transfer, err := s.paymentGateWay.CreateTransfer(int64(paymentData.Amount), recipient.RecipientID)
	if err != nil {
		return err
	}

	transactionUUID, err := uuid.Parse(paymentData.TransactionID)
	if err != nil {
		return err
	}
	logPayment := entities.LogPayment{
		TransactionUUID:    transactionUUID,
		OmiseTransactionID: transfer.ID,
	}

	if err := s.paymentRepo.CreateLogPayment(&logPayment); err != nil {
		return err
	}

	return err
}
