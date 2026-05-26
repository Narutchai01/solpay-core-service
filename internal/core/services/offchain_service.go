package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/google/uuid"
)

type OffChainService interface {
	ComFirmOffchain(ctx context.Context, req request.OffChainRequest, accountID uint) (entities.TransactionEntity, error)
}

type offChainService struct {
	transactionRepo ports.TransactionRepository
	uow             ports.UnitOfWork
	publisher       ports.Publisher
	quoteRepo       ports.QuoteRepository
}

func NewOffChainService(transactionRepo ports.TransactionRepository, uow ports.UnitOfWork, publisher ports.Publisher, quoteRepo ports.QuoteRepository) OffChainService {
	return &offChainService{
		transactionRepo: transactionRepo,
		uow:             uow,
		publisher:       publisher,
		quoteRepo:       quoteRepo,
	}
}

func (s *offChainService) ComFirmOffchain(ctx context.Context, req request.OffChainRequest, accountID uint) (entities.TransactionEntity, error) {

	cfg := config.LoadConfig()
	quoteID := strings.TrimSpace(req.QuoteID)
	if quoteID == "" {
		return entities.TransactionEntity{}, entities.NewAppError(entities.ErrTypeBadRequest, "quote_id is required", nil)
	}

	quote, err := fetchQuote(s.quoteRepo, quoteID)
	if err != nil {
		return entities.TransactionEntity{}, err
	}

	if err := validateOffchainQuote(quote, accountID); err != nil {
		return entities.TransactionEntity{}, err
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
		TransactionType: string(entities.OFFCHAIN),
		THBAmount:       float64(quote.THBAmount),
		USDTAmount:      0,
		Fee:             0,
	}

	result, err := s.uow.Execute(ctx, func(txCtx context.Context) (any, error) {
		if err := s.transactionRepo.CreateTransaction(txCtx, tx); err != nil {
			return nil, err
		}

		if quote.PromptPayID == nil || *quote.PromptPayID == "" {
			return nil, entities.NewAppError(entities.ErrTypeBadRequest, "promptpay_id is required for offchain confirm", nil)
		}

		txOffChain := &entities.TransactionOffChain{
			TransactionID: tx.TransactionUUID,
			PromptPayID:   *quote.PromptPayID,
		}

		if err := s.transactionRepo.CreateTransactionOffChain(txCtx, txOffChain); err != nil {
			return nil, err
		}
		return tx, nil
	})

	if err != nil {
		return entities.TransactionEntity{}, fmt.Errorf("create transaction: %v", err)
	}

	tx = result.(*entities.TransactionEntity)

	s.publisher.PublishTransactionMessage(cfg.TRANSACTION_ORCHESTRATOR_QUEUE, tx, string(entities.StatusBalanceWithdrawing))

	return *tx, nil
}

func validateOffchainQuote(quote *entities.Quote, accountID uint) error {
	if quote.ExpiresAt.Before(time.Now()) {
		return entities.NewAppError(entities.ErrTypeBadRequest, "quote has expired", nil)
	}

	if quote.AccountID != accountID {
		return entities.NewAppError(entities.ErrTypeConflict, "quote does not belong to the account", nil)
	}

	if quote.Status != string(entities.ACTIVE) {
		return entities.NewAppError(entities.ErrTypeBadRequest, "quote is not active", nil)
	}

	return nil
}
