package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/models"
	"github.com/google/uuid"
)

const errPublisherNotConfigured = "publisher is not configured"
const errUpdateTransactionStatusFmt = "update transaction status: %w"

// TransactionService defines operations for managing transactions.
type TransactionService interface {
	CreateTransaction(ctx context.Context, req request.CreateTransactionRequest) (*entities.TransactionEntity, error)
	GetTransactionByID(id int) (*entities.TransactionEntity, error)
	HandleTransactionUpdate(ctx context.Context, msg []byte) error
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

	if event.TxID == "" || event.Status == "" {
		return entities.NewAppError(entities.ErrTypeBadRequest, "tx_id and status are required", nil)
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
		return nil
	}

	if err := s.transactionRepo.UpdateTransactionStatus(ctx, event.TxID, event.Status); err != nil {
		return fmt.Errorf(errUpdateTransactionStatusFmt, err)
	}

	return s.routeEvent(ctx, tx, event)
}

func (s *transactionService) routeEvent(ctx context.Context, tx *entities.TransactionEntity, event request.TransactionMessage) error {
	switch tx.TransactionType {
	case string(entities.TOPUP):
		return s.topupWorkflow(ctx, tx, event)
	case string(entities.OFFCHAIN):
		return s.offchainWorkflow(ctx, tx, event)
	default:
		return nil
	}
}

func (s *transactionService) topupWorkflow(ctx context.Context, tx *entities.TransactionEntity, event request.TransactionMessage) error {
	switch event.Status {
	case string(entities.StatusSolanaSubmitted):
		if err := s.transactionRepo.UpdateTransactionStatus(ctx, tx.TransactionUUID.String(), string(entities.StatusSolanaSubmitted)); err != nil {
			return fmt.Errorf(errUpdateTransactionStatusFmt, err)
		}
		return s.publishBlockchainTransaction(tx)
	case string(entities.StatusSolanaFailed):
		if err := s.transactionRepo.UpdateTransactionStatus(ctx, tx.TransactionUUID.String(), string(entities.StatusSolanaFailed)); err != nil {
			return fmt.Errorf(errUpdateTransactionStatusFmt, err)
		}
	case string(entities.StatusSolanaSuccess):
		return s.publishBalanceTransaction(tx, string(entities.ActionDeposit))
	case string(entities.StatusBalanceUpdated):
		if err := s.transactionRepo.UpdateTransactionStatus(ctx, tx.TransactionUUID.String(), string(entities.StatusCompleted)); err != nil {
			return fmt.Errorf(errUpdateTransactionStatusFmt, err)
		}
	default:
		return fmt.Errorf("unhandled transaction status: %s", event.Status)
	}

	return nil
}

func (s *transactionService) offchainWorkflow(ctx context.Context, tx *entities.TransactionEntity, event request.TransactionMessage) error {
	switch event.Status {
	case string(entities.StatusBalanceWithdrawing):
		return s.publishBalanceTransaction(tx, string(entities.ActionWithdraw))
	case string(entities.StatusBalanceUpdating):
		return s.publishBalanceTransaction(tx, string(entities.ActionWithdraw))
	case string(entities.StatusBalanceUpdated):
		if err := s.transactionRepo.UpdateTransactionStatus(ctx, tx.TransactionUUID.String(), string(entities.StatusBalanceUpdated)); err != nil {
			return fmt.Errorf(errUpdateTransactionStatusFmt, err)
		}
		return s.publishPaymentTransaction(tx)
	case string(entities.StatusPaymentSuccess):
		if err := s.transactionRepo.UpdateTransactionStatus(ctx, tx.TransactionUUID.String(), string(entities.StatusCompleted)); err != nil {
			return fmt.Errorf(errUpdateTransactionStatusFmt, err)
		}
	default:
		return fmt.Errorf("unhandled transaction status: %s", event.Status)
	}

	return nil
}

func (s *transactionService) publishBlockchainTransaction(tx *entities.TransactionEntity) error {
	if s.publisher == nil {
		return entities.NewAppError(entities.ErrTypeInternal, errPublisherNotConfigured, nil)
	}
	if tx.TransactionOnChain == nil || tx.TransactionOnChain.TxHash == "" {
		return entities.NewAppError(entities.ErrTypeBadRequest, "missing on-chain transaction payload", nil)
	}

	cfg := config.LoadConfig()

	msg := models.SolanaTxMessage{
		TxID:     tx.TransactionUUID.String(),
		Base64Tx: tx.TransactionOnChain.TxHash,
		MetaData: models.SolanaMetaData{
			TransactionType: tx.TransactionType,
			AmountTHB:       int(tx.THBAmount),
			AmountUSDC:      int(tx.USDTAmount),
			AccountID:       int(tx.AccountID),
		},
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal solana tx message: %w", err)
	}

	return s.publisher.Publish(cfg.SOLANA_WORK_QUEUE, jsonMsg)
}

func (s *transactionService) publishBalanceTransaction(tx *entities.TransactionEntity, action string) error {
	if s.publisher == nil {
		return entities.NewAppError(entities.ErrTypeInternal, errPublisherNotConfigured, nil)
	}

	cfg := config.LoadConfig()
	msg := request.UpdateBalanceCommand{
		AccountID:     uint(tx.AccountID),
		TransactionID: tx.TransactionUUID.String(),
		THBAmount:     int64(tx.THBAmount),
		Action:        action,
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal balance message: %w", err)
	}

	err = s.publisher.Publish(cfg.BALANCE_QUEUE, jsonMsg)
	if err != nil {
		return fmt.Errorf("publish balance message: %w", err)
	}

	return nil
}

func (s *transactionService) publishPaymentTransaction(tx *entities.TransactionEntity) error {
	if s.publisher == nil {
		return entities.NewAppError(entities.ErrTypeInternal, errPublisherNotConfigured, nil)
	}
	if tx.TransactionOffChain == nil || tx.TransactionOffChain.PromptPayID == "" {
		return entities.NewAppError(entities.ErrTypeBadRequest, "missing off-chain transaction payload", nil)
	}

	cfg := config.LoadConfig()
	msg := request.RequestPaymentQueue{
		TransactionID: string(tx.TransactionUUID.String()),
		Amount:        int64(tx.THBAmount),
		Number:        tx.TransactionOffChain.PromptPayID,
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal balance message: %w", err)
	}

	err = s.publisher.Publish(cfg.PAYMENT_QUEUE, jsonMsg)
	if err != nil {
		return fmt.Errorf("publish balance message: %w", err)
	}

	return nil
}
