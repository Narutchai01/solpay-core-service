package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TopUpService interface {
	ComfirmTopUp(ctx context.Context, req request.TopUpRequest) (entities.TransactionEntity, error)
}

type topUpService struct {
	transactionRepo ports.TransactionRepository
	quoteRepo       ports.QuoteRepository
	uow             ports.UnitOfWork
}

func NewTopUpService(
	transactionRepo ports.TransactionRepository,
	quoteRepo ports.QuoteRepository,
	uow ports.UnitOfWork,
) TopUpService {
	return &topUpService{
		transactionRepo: transactionRepo,
		quoteRepo:       quoteRepo,
		uow:             uow,
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

	if quote.Status != string(entities.USED) {
		return entities.TransactionEntity{}, entities.NewAppError(entities.ErrTypeBadRequest, "quote has expired", nil)
	}

	req.SetDefaultSlippage()
	err = utils.VerifyWithSlippage(quote.ExchangeRate, 32.0, *req.MaxSlippage)
	if err != nil {
		return entities.TransactionEntity{}, entities.NewAppError(entities.ErrTypeBadRequest, "slippage too high", err)
	}

	txUUID, err := uuid.NewV7()
	if err != nil {
		return entities.TransactionEntity{}, fmt.Errorf("generate UUID: %w", err)
	}

	tx := &entities.TransactionEntity{
		TransactionUUID: txUUID,
		AccountID:       1,
		TransactionType: string(entities.TOPUP),
		THBAmount:       float64(quote.THBAmount),
		USDTAmount:      quote.USDTAmount,
		Fee:             quote.Fee,
	}

	result, err := s.uow.Execute(ctx, func(txCtx context.Context) (any, error) {
		if err := s.transactionRepo.CreateTransaction(txCtx, tx); err != nil {
			return nil, err
		}

		txOnChain := &entities.TransactionOnChain{
			TransactionID: tx.TransactionUUID,
			TxHash:        req.TxHash,
		}

		if err := s.transactionRepo.CreateTransactionOnChain(txCtx, txOnChain); err != nil {
			return nil, err
		}
		return tx, nil
	})

	if err != nil {
		return entities.TransactionEntity{}, fmt.Errorf("create transaction: %v", err)
	}

	tx = result.(*entities.TransactionEntity)

	return *tx, nil
}
