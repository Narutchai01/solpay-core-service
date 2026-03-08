package services

import (
	"context"
	"log"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
)

type PaymentService interface {
	ProcessPayment(ctx context.Context, transaction entities.TransactionEntity) error
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

func (s *paymentService) ProcessPayment(ctx context.Context, transaction entities.TransactionEntity) error {
	promptPayID := transaction.TransactionOffChain.PropmtPayID

	recipient, err := s.paymentRepo.GetRecipentByNumber(promptPayID)
	if err != nil {
		if err == entities.ErrNotFound {

			rpRaw := request.CreateRecipient{
				Name:   "Joh doe",
				Number: promptPayID,
			}

			rp, err := s.paymentGateWay.CreateRecipient(rpRaw)
			if err != nil {
				return err
			}

			recipient = entities.Recipient{Number: promptPayID, RecipientID: rp.ID}
			if err := s.paymentRepo.CreateRecipient(&recipient); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	transfer, err := s.paymentGateWay.CreateTransfer(int64(transaction.THBAmount), recipient.RecipientID)

	log.Printf("Transfer created: %+v", transfer)

	return err
}
