package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"gorm.io/gorm"
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
	if req.QuoteID == "" || req.TxHash == "" {
		return entities.TransactionEntity{}, entities.NewAppError(entities.ErrTypeBadRequest, "quote_id and tx_hash are required", nil)
	}

	quote, err := s.quoteRepo.GetQuoteByID(req.QuoteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, entities.ErrNotFound) {
			return entities.TransactionEntity{}, entities.NewAppError(entities.ErrTypeNotFound, fmt.Sprintf("quote %s not found", req.QuoteID), err)
		}
		return entities.TransactionEntity{}, entities.NewAppError(entities.ErrTypeInternal, "failed to fetch quote", err)
	}

	if time.Now().After(quote.ExpiresAt) {
		return entities.TransactionEntity{}, entities.NewAppError(entities.ErrTypeBadRequest, "quote has expired", nil)
	}

	rawTx := request.CreateTransactionRequest{
		TransactionType: string(entities.TOPUP),
		THBAmount:       float64(quote.THBAmount) / 100,
		USDTAmount:      quote.USDTAmount,
		TxHash:          &req.TxHash,
		Fee:             quote.Fee,
	}

	tx, err := s.transactionService.CreateTransaction(ctx, rawTx)
	if err != nil {
		return entities.TransactionEntity{}, entities.NewAppError(entities.ErrTypeInternal, "failed to create topup transaction", err)
	}

	return *tx, nil
}
