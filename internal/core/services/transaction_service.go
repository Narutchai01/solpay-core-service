package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/models"
	"github.com/google/uuid"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, req request.CreateTransactionRequest) (*entities.TransactionEntity, error)
	GetTransactionByID(id int) (*entities.TransactionEntity, error)
	HandleTransactionUpdate(ctx context.Context, msg []byte) error
}

type transactionService struct {
	transactionRepo ports.TransactionRepository
	uowRepo         ports.UnitOfWork
	publisher       ports.Publisher
	wsHub           ports.WebSocketPort
}

func NewTransactionService(transactionRepo ports.TransactionRepository, uowRepo ports.UnitOfWork, publisher ports.Publisher, wsHub ports.WebSocketPort) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		uowRepo:         uowRepo,
		publisher:       publisher,
		wsHub:           wsHub,
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

	switch req.TransactionType {
	case string(entities.TOPUP):
		err := s.handleMqBlockchainTransaction(result.(*entities.TransactionEntity).TransactionUUID)
		if err != nil {
			return nil, err
		}
	case string(entities.OFFCHAIN):
		err := s.handlerMqOffChainTransaction(result.(*entities.TransactionEntity).TransactionUUID)
		if err != nil {
			return nil, err
		}
	}

	return result.(*entities.TransactionEntity), nil
}

func (s *transactionService) handleCreateTransactionType(ctx context.Context, req request.CreateTransactionRequest, txId uuid.UUID) error {
	switch req.TransactionType {
	case string(entities.TOPUP):
		return s.createTransactionOnChain(ctx, req, txId)
	case string(entities.OFFCHAIN):
		return s.createTransactionOffChain(ctx, req, txId)
	case string(entities.ONCHAIN):
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

func (s *transactionService) handleMqBlockchainTransaction(txUUID uuid.UUID) error {

	cfg := config.LoadConfig()

	transaction, err := s.transactionRepo.GetTransactionByUUID(txUUID)
	if err != nil {
		return err
	}

	metaData := models.MetaDataSolana{
		Transaction_type: transaction.TransactionType,
		Amount_THB:       int(transaction.THBAmount),
		Amount_USDC:      int(transaction.USDTAmount),
		AccountID:        int(transaction.AccountID),
	}

	solanaTxMessage := models.SolanaTxMessage{
		TxID:     transaction.TransactionUUID.String(),
		Base64Tx: transaction.TransactionOnChain.TxHash,
		MetaData: metaData,
	}

	jsonMessage, err := json.Marshal(solanaTxMessage)
	if err != nil {
		return err
	}

	s.publisher.Publish(cfg.SOLANA_WORK_QUEUE, jsonMessage)

	return nil
}

func (s *transactionService) HandleTransactionUpdate(ctx context.Context, msg []byte) error {
	var txMsg request.TransactionMessage
	err := json.Unmarshal(msg, &txMsg)
	if err != nil {
		return err
	}

	txUUID, err := uuid.Parse(txMsg.TxID)
	if err != nil {
		return err
	}

	tx, err := s.transactionRepo.GetTransactionByUUID(txUUID)
	if err != nil {
		return err
	}

	if tx.Status == txMsg.Status {
		log.Printf("Transaction %s already has status %s", txMsg.TxID, txMsg.Status)
		return errors.New("transaction already processed")
	}

	err = s.transactionRepo.UpdateTransactionStatus(ctx, txMsg.TxID, txMsg.Status)
	if err != nil {
		return err
	}

	return s.processStateTransaction(ctx, tx, txMsg)
}

func (s *transactionService) processStateTransaction(ctx context.Context, tx *entities.TransactionEntity, event request.TransactionMessage) error {

	switch tx.TransactionType {
	case "top_up":
		return s.topUpWorkflow(ctx, tx, event)
	}
	return nil
}

func (s *transactionService) topUpWorkflow(ctx context.Context, tx *entities.TransactionEntity, event request.TransactionMessage) error {
	cfg := config.LoadConfig()

	switch event.SourceWorker {
	case "SOLANA":
		return s.handleSolanaWorker(ctx, tx, event, cfg)
	case "BALANCE":
		return s.handleBalanceWorker(ctx, event)
	}

	return nil
}

func (s *transactionService) handleSolanaWorker(ctx context.Context, tx *entities.TransactionEntity, event request.TransactionMessage, cfg *config.Config) error {
	switch event.Status {
	case string(entities.StatusSolanaSuccess):
		rawMsg := request.UpdateBalanceCommand{
			TransactionID: tx.TransactionUUID.String(),
			AccountID:     tx.AccountID,
			THBAmount:     int64(tx.THBAmount),
			USDTAmount:    int64(tx.USDTAmount),
		}

		jsonMessage, err := json.Marshal(rawMsg)
		if err != nil {
			return err
		}

		if err := s.transactionRepo.UpdateTransactionStatus(ctx, event.TxID, string(entities.StatusBalanceUpdating)); err != nil {
			return err
		}

		s.publisher.Publish(cfg.BALANCE_QUEUE, jsonMessage)
		s.wsHub.NotifyTransactionStatus(event.TxID, []byte(`{"status":"BALANCE_UPDATING"}`))
	case string(entities.StatusSolanaFailed):
		if err := s.transactionRepo.UpdateTransactionStatus(ctx, event.TxID, string(entities.StatusRefunded)); err != nil {
			return err
		}

		s.wsHub.NotifyTransactionStatus(event.TxID, []byte(`{"status":"FAILED"}`))
	}
	return nil
}

func (s *transactionService) handleBalanceWorker(ctx context.Context, event request.TransactionMessage) error {
	switch event.Status {
	case string(entities.StatusBalanceUpdated):
		if err := s.transactionRepo.UpdateTransactionStatus(ctx, event.TxID, string(entities.StatusCompleted)); err != nil {
			return err
		}
		s.wsHub.NotifyTransactionStatus(event.TxID, []byte(`{"status":"COMPLETED"}`))
	case string(entities.StatusBalanceFailed):
		if err := s.transactionRepo.UpdateTransactionStatus(ctx, event.TxID, string(entities.StatusRefunded)); err != nil {
			return err
		}
		s.wsHub.NotifyTransactionStatus(event.TxID, []byte(`{"status":"REFUNDED"}`))
	}
	return nil
}

func (s *transactionService) handlerMqOffChainTransaction(txUUID uuid.UUID) error {
	cfg := config.LoadConfig()

	transaction, err := s.transactionRepo.GetTransactionByUUID(txUUID)
	if err != nil {
		return err
	}

	promptPayMessage := request.RequestPaymentQueue{
		TransactionID: transaction.TransactionUUID.String(),
		Number:        transaction.TransactionOffChain.PropmtPayID,
		Amount:        int64(transaction.THBAmount),
	}

	jsonMessage, err := json.Marshal(promptPayMessage)
	if err != nil {
		return err
	}

	s.publisher.Publish(cfg.PAYMENT_QUEUE, jsonMessage)

	return nil
}
