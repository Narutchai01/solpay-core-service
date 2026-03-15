package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/models"
	"github.com/google/uuid"
)

// TransactionService defines operations for managing transactions.
type TransactionService interface {
	CreateTransaction(ctx context.Context, req request.CreateTransactionRequest) (*entities.TransactionEntity, error)
	GetTransactionByID(id int) (*entities.TransactionEntity, error)
	HandleTransactionUpdate(ctx context.Context, msg []byte) error
	CreateTransactionTopUp(ctx context.Context, req models.CreateTransactionTopUp) (*entities.TransactionEntity, error)
}

type transactionService struct {
	transactionRepo ports.TransactionRepository
	uow             ports.UnitOfWork
	publisher       ports.Publisher
	wsHub           ports.WebSocketPort
}

// NewTransactionService creates a new TransactionService.
func NewTransactionService(
	transactionRepo ports.TransactionRepository,
	uow ports.UnitOfWork,
	publisher ports.Publisher,
	wsHub ports.WebSocketPort,
) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		uow:             uow,
		publisher:       publisher,
		wsHub:           wsHub,
	}
}

// ---------------------------------------------------------------------------
// CreateTransaction
// ---------------------------------------------------------------------------

func (s *transactionService) CreateTransaction(ctx context.Context, req request.CreateTransactionRequest) (*entities.TransactionEntity, error) {
	txUUID, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generate UUID: %w", err)
	}

	transaction := &entities.TransactionEntity{
		TransactionUUID: txUUID,
		AccountID:       1,
		TransactionType: req.TransactionType,
		THBAmount:       req.THBAmount,
		USDTAmount:      req.USDTAmount,
		Fee:             req.Fee,
	}

	result, err := s.uow.Execute(ctx, func(txCtx context.Context) (any, error) {
		if err := s.transactionRepo.CreateTransaction(txCtx, transaction); err != nil {
			return nil, err
		}
		if err := s.createSubTransactions(txCtx, req, txUUID); err != nil {
			return nil, err
		}
		return transaction, nil
	})
	if err != nil {
		return nil, fmt.Errorf("create transaction: %w", err)
	}

	tx := result.(*entities.TransactionEntity)

	return tx, nil
}

// createSubTransactions creates the type-specific child records within the UoW.
func (s *transactionService) createSubTransactions(ctx context.Context, req request.CreateTransactionRequest, txID uuid.UUID) error {
	switch req.TransactionType {
	case string(entities.TOPUP):
		return s.createTransactionOnChain(ctx, req, txID)
	case string(entities.OFFCHAIN):
		return s.createTransactionOffChain(ctx, req, txID)
	case string(entities.ONCHAIN):
		if err := s.createTransactionOnChain(ctx, req, txID); err != nil {
			return err
		}
		return s.createTransactionOffChain(ctx, req, txID)
	default:
		return fmt.Errorf("unsupported transaction type: %s", req.TransactionType)
	}
}

func (s *transactionService) createTransactionOnChain(ctx context.Context, req request.CreateTransactionRequest, txID uuid.UUID) error {
	if req.TxHash == nil {
		return errors.New("tx_hash is required for on-chain transactions")
	}
	return s.transactionRepo.CreateTransactionOnChain(ctx, &entities.TransactionOnChain{
		TransactionID: txID,
		TxHash:        *req.TxHash,
	})
}

func (s *transactionService) createTransactionOffChain(ctx context.Context, req request.CreateTransactionRequest, txID uuid.UUID) error {
	if req.PromptPayID == nil {
		return errors.New("prompt_pay_id is required for off-chain transactions")
	}
	return s.transactionRepo.CreateTransactionOffChain(ctx, &entities.TransactionOffChain{
		TransactionID: txID,
		PromptPayID:   *req.PromptPayID,
	})
}

// ---------------------------------------------------------------------------
// GetTransactionByID
// ---------------------------------------------------------------------------

func (s *transactionService) GetTransactionByID(id int) (*entities.TransactionEntity, error) {
	tx, err := s.transactionRepo.GetTransactionByID(id)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			return nil, entities.NewAppError(entities.ErrTypeNotFound, fmt.Sprintf("transaction %d not found", id), err)
		}
		return nil, entities.NewAppError(entities.ErrTypeInternal, "failed to get transaction", err)
	}
	return tx, nil
}

func (s *transactionService) HandleTransactionUpdate(ctx context.Context, msg []byte) error {
	var event request.TransactionMessage
	if err := json.Unmarshal(msg, &event); err != nil {
		return fmt.Errorf("unmarshal transaction message: %w", err)
	}

	txUUID, err := uuid.Parse(event.TxID)
	if err != nil {
		return fmt.Errorf("parse transaction UUID: %w", err)
	}

	tx, err := s.transactionRepo.GetTransactionByUUID(txUUID)
	if err != nil {
		return fmt.Errorf("get transaction by UUID: %w", err)
	}

	if tx.Status == event.Status {
		log.Printf("Transaction %s already has status %s, skipping", event.TxID, event.Status)
		return errors.New("transaction already processed")
	}

	if err := s.transactionRepo.UpdateTransactionStatus(ctx, event.TxID, event.Status); err != nil {
		return fmt.Errorf("update transaction status: %w", err)
	}

	return nil
}

func (s *transactionService) CreateTransactionTopUp(ctx context.Context, req models.CreateTransactionTopUp) (*entities.TransactionEntity, error) {

	return nil, nil
}
