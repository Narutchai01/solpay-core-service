package services

import (
	"context"
	"errors"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports/repositories"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/google/uuid"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, req request.CreateTransactionRequest) (*entities.TransactionEntity, error)
	GetTransactionByID(id int) (*entities.TransactionEntity, error)
}

type transactionService struct {
	transactionRepo repositories.TransactionRepository
	uowRepo         repositories.UnitOfWork
}

func NewTransactionService(transactionRepo repositories.TransactionRepository, uowRepo repositories.UnitOfWork) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		uowRepo:         uowRepo,
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

	result, err := s.uowRepo.Execute(ctx, func(ctx context.Context) (any, error) {
		if err := s.transactionRepo.CreateTransaction(ctx, transaction); err != nil {
			return nil, err
		}

		err := s.handleCreateTransactionType(ctx, req, genreateUUID)
		if err != nil {
			return nil, err
		}

		return transaction, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*entities.TransactionEntity), nil
}

func (s *transactionService) handleCreateTransactionType(ctx context.Context, req request.CreateTransactionRequest, txId uuid.UUID) error {
	switch req.TransactionType {
	case "top_up":
		return s.createTransactionOnChain(ctx, req, txId)
	case "transaction_offchain":
		return s.createTransactionOffChain(ctx, req, txId)
	case "transaction_onchain":
		err := s.createTransactionOnChain(ctx, req, txId)
		if err != nil {
			return err
		}
		return s.createTransactionOffChain(ctx, req, txId)
	default:
		return errors.New("invalid transaction type")
	}
}

func (s *transactionService) createTransactionOnChain(ctx context.Context, req request.CreateTransactionRequest, txId uuid.UUID) error {
	if req.TxHash == nil {
		return errors.New("onchain require")
	}
	transactionOnChain := &entities.TransactionOnChain{
		TransactionID: txId,
		TxHash:        *req.TxHash,
	}
	err := s.transactionRepo.CreateTransactionOnChain(ctx, transactionOnChain)
	if err != nil {
		return err
	}
	return nil
}

func (s *transactionService) createTransactionOffChain(ctx context.Context, req request.CreateTransactionRequest, txId uuid.UUID) error {
	if req.PromptPayID == nil {
		return errors.New("offchain require")
	}
	transactionOffChain := &entities.TransactionOffChain{
		TransactionID: txId,
		PropmtPayID:   *req.PromptPayID,
	}
	err := s.transactionRepo.CreateTransactionOffChain(ctx, transactionOffChain)
	if err != nil {
		return err
	}
	return nil
}

func (s *transactionService) GetTransactionByID(id int) (*entities.TransactionEntity, error) {
	transaction, err := s.transactionRepo.GetTransactionByID(id)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			return &entities.TransactionEntity{}, entities.NewAppError(entities.ErrTypeNotFound, "Transaction not found", err)
		}
		return &entities.TransactionEntity{}, entities.NewAppError(entities.ErrTypeInternal, "Failed to get transaction", err)
	}
	return transaction, nil
}
