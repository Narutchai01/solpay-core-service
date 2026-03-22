package services

import (
	"context"
	"fmt"

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
}

func NewOffChainService(transactionRepo ports.TransactionRepository, uow ports.UnitOfWork, publisher ports.Publisher) OffChainService {
	return &offChainService{
		transactionRepo: transactionRepo,
		uow:             uow,
		publisher:       publisher,
	}
}

func (s *offChainService) ComFirmOffchain(ctx context.Context, req request.OffChainRequest, accountID uint) (entities.TransactionEntity, error) {

	cfg := config.LoadConfig()

	txUUID, err := uuid.NewV7()
	if err != nil {
		return entities.TransactionEntity{}, fmt.Errorf("generate UUID: %w", err)
	}

	tx := &entities.TransactionEntity{
		TransactionUUID: txUUID,
		AccountID:       accountID,
		TransactionType: string(entities.OFFCHAIN),
		THBAmount:       float64(req.THBAmount),
		USDTAmount:      0,
		Fee:             0,
	}

	result, err := s.uow.Execute(ctx, func(txCtx context.Context) (any, error) {
		if err := s.transactionRepo.CreateTransaction(txCtx, tx); err != nil {
			return nil, err
		}

		txOffChain := &entities.TransactionOffChain{
			TransactionID: tx.TransactionUUID,
			PromptPayID:   req.PromptPayID,
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
