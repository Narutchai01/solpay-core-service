package services

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
)

type BalanceService interface {
	GetBalances(page int, limit int) ([]entities.BalanceEntity, int64, error)
	GetBalanceByID(id int) (*entities.BalanceEntity, error)
	UpdateBalance(data []byte) error
}

type balanceService struct {
	accountRepo ports.AccountRepository
	balanceRepo ports.BalanceRepository
	uowRepo     ports.UnitOfWork
}

func NewBalanceService(balanceRepo ports.BalanceRepository, uowRepo ports.UnitOfWork) BalanceService {
	return &balanceService{
		balanceRepo: balanceRepo,
		uowRepo:     uowRepo,
	}
}

func (s *balanceService) GetBalances(page int, limit int) ([]entities.BalanceEntity, int64, error) {
	var balances []entities.BalanceEntity
	var total int64
	var errList, errCount error
	var wg sync.WaitGroup

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

	if errList != nil {
		msg := utils.FormatValidationError(errList)
		return []entities.BalanceEntity{}, 0, entities.NewAppError(entities.ErrTypeInternal, msg, errList)
	}

	if errCount != nil {
		msg := utils.FormatValidationError(errCount)
		return []entities.BalanceEntity{}, 0, entities.NewAppError(entities.ErrTypeInternal, msg, errCount)
	}

	return balances, total, nil
}

func (s *balanceService) GetBalanceByID(id int) (*entities.BalanceEntity, error) {
	balance, err := s.balanceRepo.GetBalanceByID(id)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			return &entities.BalanceEntity{}, entities.NewAppError(entities.ErrTypeNotFound, "balance not found", err)
		}
		return &entities.BalanceEntity{}, entities.NewAppError(entities.ErrTypeInternal, "internal server error", err)
	}
	return balance, nil
}

func (s *balanceService) UpdateBalance(data []byte) error {
	var cmd request.UpdateBalanceCommand
	err := json.Unmarshal(data, &cmd)
	if err != nil {
		return entities.NewAppError(entities.ErrTypeBadRequest, "invalid request body", err)
	}

	newBalance := &entities.BalanceEntity{
		AccountID:  cmd.AccountID,
		THBAmount:  cmd.THBAmount,
		USDTAmount: cmd.USDTAmount,
	}

	err = s.balanceRepo.UpdateBalance(context.Background(), newBalance)
	if err != nil {
		return entities.NewAppError(entities.ErrTypeInternal, "failed to update balance", err)
	}

	return nil
}
