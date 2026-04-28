package services

import (
	"context"
	"fmt"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/google/uuid"
)

type TopUpService interface {
	ComfirmTopUp(ctx context.Context, req request.TopUpRequest, accountID uint) (entities.TransactionEntity, error)
}

type topUpService struct {
	transactionRepo ports.TransactionRepository
	quoteRepo       ports.QuoteRepository
	uow             ports.UnitOfWork
	publisher       ports.Publisher
}

func NewTopUpService(
	transactionRepo ports.TransactionRepository,
	quoteRepo ports.QuoteRepository,
	uow ports.UnitOfWork,
	publisher ports.Publisher,
) TopUpService {
	return &topUpService{
		transactionRepo: transactionRepo,
		quoteRepo:       quoteRepo,
		uow:             uow,
		publisher:       publisher,
	}
}

func (s *topUpService) ComfirmTopUp(ctx context.Context, req request.TopUpRequest, accountID uint) (entities.TransactionEntity, error) {
	cfg := config.LoadConfig()

	if err := validateTopUpRequest(req); err != nil {
		return entities.TransactionEntity{}, err
	}

	quote, err := fetchQuote(s.quoteRepo, req.QuoteID)
	if err != nil {
		return entities.TransactionEntity{}, err
	}

	if err := validateQuote(quote); err != nil {
		return entities.TransactionEntity{}, err
	}

	req.SetDefaultSlippage()
	if err := utils.VerifyWithSlippage(quote.ExchangeRate, 32.0, *req.MaxSlippage); err != nil {
		return entities.TransactionEntity{}, entities.NewAppError(entities.ErrTypeBadRequest, "slippage too high", err)
	}

	txUUID, err := uuid.NewV7()
	if err != nil {
		return entities.TransactionEntity{}, fmt.Errorf("generate UUID: %w", err)
	}

	categoryID := req.CategoryID
	if categoryID == 0 {
		categoryID = 1
	}

	tx := &entities.TransactionEntity{
		TransactionUUID: txUUID,
		AccountID:       accountID,
		CategoryID:      categoryID,
		TransactionType: string(entities.TOPUP),
		THBAmount:       float64(quote.THBAmount),
		USDTAmount:      quote.USDTAmount,
		Fee:             quote.Fee,
	}

	result, err := s.uow.Execute(ctx, func(txCtx context.Context) (any, error) {
		if err := s.transactionRepo.CreateTransaction(txCtx, tx); err != nil {
			return nil, err
		}

		signature := signatureFromTxHash(req.TxHash)

		txOnChain := &entities.TransactionOnChain{
			TransactionID: tx.TransactionUUID,
			Signature:     signature,
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

	s.publisher.PublishTransactionMessage(cfg.TRANSACTION_ORCHESTRATOR_QUEUE, tx, string(entities.StatusSolanaSubmitted))

	return *tx, nil
}
