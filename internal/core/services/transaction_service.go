package services

import (
	"context"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports/repositories"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
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
	transaction := &entities.TransactionEntity{
		AccountID:  req.AccountID,
		CategoryID: req.CategoryID,
		Type:       req.Type,
		Status:     req.Status,
		THBAmount:  req.THBAmount,
		USDTAmount: req.USDTAmount,
		Fee:        req.Fee,
	}

	if err := s.transactionRepo.CreateTransaction(ctx, transaction); err != nil {
		return &entities.TransactionEntity{}, err
	}

	return transaction, nil
}
