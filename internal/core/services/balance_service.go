package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
)

// BalanceService defines operations for managing balances.
type BalanceService interface {
	GetBalances(page, limit int) ([]entities.BalanceEntity, int64, error)
	GetBalanceByID(id int) (*entities.BalanceEntity, error)
	UpdateBalance(data []byte) error
	WithDraw(data []byte) error
	Deposit(data []byte) error
	GetByBalanceByAccountID(accountID uint) (*entities.BalanceEntity, error)
}

type balanceService struct {
	balanceRepo ports.BalanceRepository
	uow         ports.UnitOfWork
	publisher   ports.Publisher
}

// NewBalanceService creates a new BalanceService.
func NewBalanceService(balanceRepo ports.BalanceRepository, uow ports.UnitOfWork, publisher ports.Publisher) BalanceService {
	return &balanceService{
		balanceRepo: balanceRepo,
		uow:         uow,
		publisher:   publisher,
	}
}

func (s *balanceService) GetBalances(page, limit int) ([]entities.BalanceEntity, int64, error) {
	var (
		balances []entities.BalanceEntity
		total    int64
		errList  error
		errCount error
		wg       sync.WaitGroup
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		balances, errList = s.balanceRepo.GetBalances(page, limit)
	}()
	go func() {
		defer wg.Done()
		total, errCount = s.balanceRepo.CountBalances()
	}()
	wg.Wait()

	if err := errors.Join(errList, errCount); err != nil {
		return nil, 0, entities.NewAppError(entities.ErrTypeInternal, "failed to list balances", err)
	}

	return balances, total, nil
}

func (s *balanceService) GetBalanceByID(id int) (*entities.BalanceEntity, error) {
	balance, err := s.balanceRepo.GetBalanceByID(id)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			return nil, entities.NewAppError(entities.ErrTypeNotFound, fmt.Sprintf("balance %d not found", id), err)
		}
		return nil, entities.NewAppError(entities.ErrTypeInternal, "failed to get balance", err)
	}
	return balance, nil
}

func (s *balanceService) UpdateBalance(data []byte) error {
	var cmd request.UpdateBalanceCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return entities.NewAppError(entities.ErrTypeBadRequest, "invalid request body", err)
	}

	switch cmd.Action {
	case string(entities.ActionDeposit):
		return s.Deposit(data)
	case string(entities.ActionWithdraw):
		return s.WithDraw(data)
	default:
		return entities.NewAppError(entities.ErrTypeBadRequest, "invalid action", nil)
	}
}

// Deposit handles deposit to a balance
func (s *balanceService) Deposit(data []byte) error {
	cfg := config.LoadConfig()

	var cmd request.UpdateBalanceCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return entities.NewAppError(entities.ErrTypeBadRequest, "invalid request body", err)
	}

	// Retrieve current balance
	balance, err := s.balanceRepo.GetBalanceByID(int(cmd.AccountID))
	if err != nil {
		s.publishBalanceResult(cfg, cmd.TransactionID, string(entities.StatusBalanceFailed))
		return entities.NewAppError(entities.ErrTypeInternal, "failed to get balance", err)
	}

	// Add amounts
	balance.THBAmount += cmd.THBAmount

	if err := s.balanceRepo.UpdateBalance(context.Background(), balance); err != nil {
		s.publishBalanceResult(cfg, cmd.TransactionID, string(entities.StatusBalanceFailed))
		return entities.NewAppError(entities.ErrTypeInternal, "failed to deposit balance", err)
	}

	s.publishBalanceResult(cfg, cmd.TransactionID, string(entities.StatusBalanceUpdated))
	return nil
}

// WithDraw handles withdrawal from a balance
func (s *balanceService) WithDraw(data []byte) error {
	cfg := config.LoadConfig()

	var cmd request.UpdateBalanceCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return entities.NewAppError(entities.ErrTypeBadRequest, "invalid request body", err)
	}

	// Retrieve current balance
	balance, err := s.balanceRepo.GetBalanceByID(int(cmd.AccountID))
	if err != nil {
		s.publishBalanceResult(cfg, cmd.TransactionID, string(entities.StatusBalanceFailed))
		return entities.NewAppError(entities.ErrTypeInternal, "failed to get balance", err)
	}

	// Check if sufficient funds
	if balance.THBAmount < cmd.THBAmount {
		s.publishBalanceResult(cfg, cmd.TransactionID, string(entities.StatusBalanceFailed))
		return entities.NewAppError(entities.ErrTypeBadRequest, "insufficient funds", nil)
	}

	// Deduct amounts
	balance.THBAmount -= cmd.THBAmount

	if err := s.balanceRepo.UpdateBalance(context.Background(), balance); err != nil {
		s.publishBalanceResult(cfg, cmd.TransactionID, string(entities.StatusBalanceFailed))
		return entities.NewAppError(entities.ErrTypeInternal, "failed to withdraw balance", err)
	}

	s.publishBalanceResult(cfg, cmd.TransactionID, string(entities.StatusBalanceUpdated))
	return nil
}

func (s *balanceService) publishBalanceResult(cfg *config.Config, txID, status string) {
	if s.publisher == nil {
		log.Printf("Publisher not configured, skipping balance result publish")
		return
	}

	msg := request.TransactionMessage{
		TxID:         txID,
		SourceWorker: "BALANCE",
		Status:       status,
	}

	jsonMessage, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal balance result message: %v", err)
		return
	}

	if err := s.publisher.Publish(cfg.TRANSACTION_ORCHESTRATOR_QUEUE, jsonMessage); err != nil {
		log.Printf("Failed to publish balance result: %v", err)
	}
}

func (s *balanceService) GetByBalanceByAccountID(accountID uint) (*entities.BalanceEntity, error) {
	balance, err := s.balanceRepo.GetBalanceByAccountID(accountID)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			return nil, entities.NewAppError(entities.ErrTypeNotFound, fmt.Sprintf("balance for account %d not found", accountID), err)
		}
		return nil, entities.NewAppError(entities.ErrTypeInternal, "failed to get balance by account ID", err)
	}
	return balance, nil
}
