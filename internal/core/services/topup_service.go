package services

import (
	"context"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
)

type TopUpService interface {
	ComfirmTopUp(ctx context.Context, req request.TopUpRequest) (entities.TransactionEntity, error)
}

type topUpService struct {
	transactionService TransactionService
	transactionRepo    ports.TransactionRepository
	quoteRepo          ports.QuoteRepository
	uow                ports.UnitOfWork
}

func NewTopUpService(
	transactionService TransactionService,
	transactionRepo ports.TransactionRepository,
	quoteRepo ports.QuoteRepository,
	uow ports.UnitOfWork,
) TopUpService {
	return &topUpService{
		transactionService: transactionService,
		transactionRepo:    transactionRepo,
		quoteRepo:          quoteRepo,
		uow:                uow,
	}
}

func (s *topUpService) ComfirmTopUp(ctx context.Context, req request.TopUpRequest) (entities.TransactionEntity, error) {

	quote, err := s.quoteRepo.GetQuoteByID(req.QuoteID)
	if err != nil {
		return entities.TransactionEntity{}, err
	}

	rawTx := request.CreateTransactionRequest{
		TransactionType: string(entities.TOPUP),
		THBAmount:       float64(quote.THBAmount),
		USDTAmount:      quote.USDTAmount,
		TxHash:          &req.TxHash,
		Fee:             quote.Fee,
	}

	tx, err := s.transactionService.CreateTransaction(ctx, rawTx)
	if err != nil {
		return entities.TransactionEntity{}, err
	}

	return *tx, nil
}
