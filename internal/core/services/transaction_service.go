package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/Narutchai01/solpay-core-service/internal/config"
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

	if err := s.publishAfterCreate(tx); err != nil {
		return nil, err
	}

	return tx, nil
}

// publishAfterCreate sends the initial MQ message for the appropriate transaction type.
func (s *transactionService) publishAfterCreate(tx *entities.TransactionEntity) error {
	switch tx.TransactionType {
	case string(entities.TOPUP):
		return s.publishBlockchainTransaction(tx.TransactionUUID)
	case string(entities.OFFCHAIN):
		return s.publishOffChainTransaction(tx.TransactionUUID)
	default:
		return nil
	}
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

// ---------------------------------------------------------------------------
// HandleTransactionUpdate (orchestrator)
// ---------------------------------------------------------------------------

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

	return s.routeEvent(ctx, tx, event)
}

// routeEvent dispatches to the correct workflow based on transaction type.
func (s *transactionService) routeEvent(ctx context.Context, tx *entities.TransactionEntity, event request.TransactionMessage) error {
	switch tx.TransactionType {
	case string(entities.TOPUP):
		return s.topUpWorkflow(ctx, tx, event)
	case string(entities.OFFCHAIN):
		return s.offchainWorkflow(ctx, event)
	default:
		return nil
	}
}

// ---------------------------------------------------------------------------
// Top-up workflow
// ---------------------------------------------------------------------------

func (s *transactionService) topUpWorkflow(ctx context.Context, tx *entities.TransactionEntity, event request.TransactionMessage) error {
	switch event.SourceWorker {
	case "SOLANA":
		return s.handleSolanaWorker(ctx, tx, event)
	case "BALANCE":
		return s.handleBalanceWorker(ctx, event)
	default:
		return nil
	}
}

func (s *transactionService) handleSolanaWorker(ctx context.Context, tx *entities.TransactionEntity, event request.TransactionMessage) error {
	cfg := config.LoadConfig()

	switch event.Status {
	case string(entities.StatusSolanaSuccess):
		cmd := request.UpdateBalanceCommand{
			TransactionID: tx.TransactionUUID.String(),
			AccountID:     tx.AccountID,
			THBAmount:     int64(tx.THBAmount),
			USDTAmount:    int64(tx.USDTAmount),
		}

		jsonMsg, err := json.Marshal(cmd)
		if err != nil {
			return fmt.Errorf("marshal balance command: %w", err)
		}

		if err := s.transactionRepo.UpdateTransactionStatus(ctx, event.TxID, string(entities.StatusBalanceUpdating)); err != nil {
			return fmt.Errorf("update status to BALANCE_UPDATING: %w", err)
		}

		if err := s.publisher.Publish(cfg.BALANCE_QUEUE, jsonMsg); err != nil {
			return fmt.Errorf("publish balance command: %w", err)
		}
		s.notifyStatus(event.TxID, "BALANCE_UPDATING")

	case string(entities.StatusSolanaFailed):
		if err := s.transactionRepo.UpdateTransactionStatus(ctx, event.TxID, string(entities.StatusRefunded)); err != nil {
			return fmt.Errorf("update status to REFUNDED: %w", err)
		}
		s.notifyStatus(event.TxID, "FAILED")
	}

	return nil
}

func (s *transactionService) handleBalanceWorker(ctx context.Context, event request.TransactionMessage) error {
	switch event.Status {
	case string(entities.StatusBalanceUpdated):
		if err := s.transactionRepo.UpdateTransactionStatus(ctx, event.TxID, string(entities.StatusCompleted)); err != nil {
			return fmt.Errorf("update status to COMPLETED: %w", err)
		}
		s.notifyStatus(event.TxID, "COMPLETED")

	case string(entities.StatusBalanceFailed):
		if err := s.transactionRepo.UpdateTransactionStatus(ctx, event.TxID, string(entities.StatusRefunded)); err != nil {
			return fmt.Errorf("update status to REFUNDED: %w", err)
		}
		s.notifyStatus(event.TxID, "REFUNDED")
	}

	return nil
}

// ---------------------------------------------------------------------------
// Off-chain workflow
// ---------------------------------------------------------------------------

func (s *transactionService) offchainWorkflow(ctx context.Context, event request.TransactionMessage) error {
	if event.SourceWorker == "PAYMENT" {
		return s.handlePaymentWorker(ctx, event)
	}
	return nil
}

func (s *transactionService) handlePaymentWorker(ctx context.Context, event request.TransactionMessage) error {
	switch event.Status {
	case string(entities.StatusPaymentSuccess):
		if err := s.transactionRepo.UpdateTransactionStatus(ctx, event.TxID, string(entities.StatusCompleted)); err != nil {
			return fmt.Errorf("update status to COMPLETED: %w", err)
		}
		s.notifyStatus(event.TxID, "COMPLETED")

	case string(entities.StatusPaymentFailed):
		if err := s.transactionRepo.UpdateTransactionStatus(ctx, event.TxID, string(entities.StatusRefunded)); err != nil {
			return fmt.Errorf("update status to REFUNDED: %w", err)
		}
		s.notifyStatus(event.TxID, "FAILED")
	}

	return nil
}

// ---------------------------------------------------------------------------
// MQ publishing helpers
// ---------------------------------------------------------------------------

func (s *transactionService) publishBlockchainTransaction(txUUID uuid.UUID) error {
	cfg := config.LoadConfig()

	tx, err := s.transactionRepo.GetTransactionByUUID(txUUID)
	if err != nil {
		return fmt.Errorf("get transaction for blockchain publish: %w", err)
	}

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

func (s *transactionService) publishOffChainTransaction(txUUID uuid.UUID) error {
	cfg := config.LoadConfig()

	tx, err := s.transactionRepo.GetTransactionByUUID(txUUID)
	if err != nil {
		return fmt.Errorf("get transaction for off-chain publish: %w", err)
	}

	msg := request.RequestPaymentQueue{
		TransactionID: tx.TransactionUUID.String(),
		Number:        tx.TransactionOffChain.PromptPayID,
		Amount:        int64(tx.THBAmount),
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal payment queue message: %w", err)
	}

	return s.publisher.Publish(cfg.PAYMENT_QUEUE, jsonMsg)
}

// notifyStatus sends a WebSocket notification if a hub is configured.
func (s *transactionService) notifyStatus(txID, status string) {
	if s.wsHub == nil {
		return
	}
	payload := fmt.Sprintf(`{"status":%q}`, status)
	s.wsHub.NotifyTransactionStatus(txID, []byte(payload))
}
