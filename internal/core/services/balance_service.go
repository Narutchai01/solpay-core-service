package services

import (
	"sync"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports/repositories"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
)

type BalanceService interface {
	GetBalances(page int, limit int) ([]entities.BalanceEntity, int64, error)
}

type balanceService struct {
	accountRepo repositories.AccountRepository
	balanceRepo repositories.BalanceRepository
	uowRepo     repositories.UnitOfWork
}

func NewBalanceService(balanceRepo repositories.BalanceRepository, uowRepo repositories.UnitOfWork) BalanceService {
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
