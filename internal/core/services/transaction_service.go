package services

import (
	"context"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports/repositories"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/google/uuid"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, req request.CreateTransactionRequest) (*entities.TransactionEntity, error)
}

type transactionService struct {
	transactionRepo repositories.TransactionRepository
}

func NewTransactionService(transactionRepo repositories.TransactionRepository) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
	}
}

func (s *transactionService) CreateTransaction(ctx context.Context, req request.CreateTransactionRequest) (*entities.TransactionEntity, error) {
	genreateUUID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	transaction := &entities.TransactionEntity{
		TransactionUUID: genreateUUID,
		AccountID:       1,
		TransactionType: req.TransactionType,
		THBAmount:       req.THBAmount,
		USDTAmount:      req.USDTAmount,
		Fee:             req.Fee,
	}
	var txOnChain *entities.TransactionOnChain
	var txOffChain *entities.TransactionOffChain

	if req.TransactionType == "top_up" || req.TransactionType == "transaction_onchain" {
		txOnChain = &entities.TransactionOnChain{
			TransactionID: genreateUUID,
			FromAddress:   *req.FromAddress,
			TxHash:        *req.TxHash,
		}
	}

	if req.TransactionType == "transaction_offchain" || req.TransactionType == "transaction_onchain" {
		txOffChain = &entities.TransactionOffChain{
			TransactionID: genreateUUID,
			PropmtPayID:   *req.PromptPayID,
		}
	}

	err = s.transactionRepo.CreateTransaction(ctx, transaction, txOnChain, txOffChain)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
