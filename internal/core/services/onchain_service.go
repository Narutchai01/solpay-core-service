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

type OnchainService interface {
	ComfirmOnchain(ctx context.Context, req request.TopUpRequest, accountID uint) (entities.TransactionEntity, error)
}

type OnchainServiceService struct {
	transactionRepo ports.TransactionRepository
	quoteRepo       ports.QuoteRepository
	uow             ports.UnitOfWork
	publisher       ports.Publisher
}

func NewOnchainService(
	transactionRepo ports.TransactionRepository,
	quoteRepo ports.QuoteRepository,
	uow ports.UnitOfWork,
	publisher ports.Publisher,
) OnchainService {
	return &OnchainServiceService{
		transactionRepo: transactionRepo,
		quoteRepo:       quoteRepo,
		uow:             uow,
		publisher:       publisher,
	}
}

func (s *OnchainServiceService) ComfirmOnchain(ctx context.Context, req request.TopUpRequest, accountID uint) (entities.TransactionEntity, error) {
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

	fmt.Printf("CategoryID in service: %d\n", req.CategoryID)

	categoryID := req.CategoryID
	if categoryID == 0 {
		categoryID = 1
	}

	tx := &entities.TransactionEntity{
		TransactionUUID: txUUID,
		AccountID:       accountID,
		TransactionType: quote.Type,
		THBAmount:       float64(quote.THBAmount),
		CategoryID:      categoryID,
		USDTAmount:      quote.USDTAmount,
		Fee:             quote.Fee,
	}
	slipURL := "https://images.unsplash.com/photo-1776320644111-f72194d35eb8?q=80&w=687&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"

	result, err := s.uow.Execute(ctx, func(txCtx context.Context) (any, error) {
		if err := s.transactionRepo.CreateTransaction(txCtx, tx); err != nil {
			return nil, err
		}

		if quote.PromptPayID == nil || *quote.PromptPayID == "" {
			return nil, entities.NewAppError(entities.ErrTypeBadRequest, "promptpay_id is required for onchain confirm", nil)
		}

		signature := signatureFromTxHash(req.TxHash)

		txOnChain := &entities.TransactionOnChain{
			TransactionID: tx.TransactionUUID,
			Signature:     signature,
			TxHash:        req.TxHash,
		}

		txOffchain := &entities.TransactionOffChain{
			TransactionID: tx.TransactionUUID,
			PromptPayID:   *quote.PromptPayID,
			SlipURL:       &slipURL,
		}

		if err := s.transactionRepo.CreateTransactionOnChain(txCtx, txOnChain); err != nil {
			return nil, err
		}

		if err := s.transactionRepo.CreateTransactionOffChain(txCtx, txOffchain); err != nil {
			return nil, err
		}
		return tx, nil
	})

	if err != nil {
		return entities.TransactionEntity{}, fmt.Errorf("create transaction: %v", err)
	}

	tx = result.(*entities.TransactionEntity)
	loadedTx, err := s.transactionRepo.GetTransactionByUUID(tx.TransactionUUID)
	if err != nil {
		return entities.TransactionEntity{}, fmt.Errorf("get transaction by uuid: %w", err)
	}
	tx = loadedTx

	s.publisher.PublishTransactionMessage(cfg.TRANSACTION_ORCHESTRATOR_QUEUE, tx, string(entities.StatusSolanaSubmitted))

	return *tx, nil
}
