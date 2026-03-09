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
	cfg := config.LoadConfig()

	var cmd request.UpdateBalanceCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return entities.NewAppError(entities.ErrTypeBadRequest, "invalid request body", err)
	}

	newBalance := &entities.BalanceEntity{
		AccountID:  cmd.AccountID,
		THBAmount:  cmd.THBAmount,
		USDTAmount: cmd.USDTAmount,
	}

	if err := s.balanceRepo.UpdateBalance(context.Background(), newBalance); err != nil {
		s.publishBalanceResult(cfg, cmd.TransactionID, string(entities.StatusBalanceFailed))
		return entities.NewAppError(entities.ErrTypeInternal, "failed to update balance", err)
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
