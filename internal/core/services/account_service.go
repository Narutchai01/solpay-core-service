package services

import (
	"context"
	"errors"
	"sync"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports/repositories"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
)

// Note: Define the AccountService interface
type AccountService interface {
	CreateAccount(ctx context.Context, req request.CreateAccountRequest) (*entities.AccountEntity, error)
	GetAccounts(page int, limit int) ([]entities.AccountEntity, int64, error)
	GetAccountByID(id int) (*entities.AccountEntity, error)
}

// Note: Implement the AccountService interface
type accountService struct {
	accountRepo repositories.AccountRepository
	balanceRepo repositories.BalanceRepository
	uowRepo     repositories.UnitOfWork
}

// Note: Constructor function for AccountService
func NewAccountService(accountRepo repositories.AccountRepository, balanceRepo repositories.BalanceRepository, uowRepo repositories.UnitOfWork) AccountService {
	return &accountService{
		accountRepo: accountRepo,
		balanceRepo: balanceRepo,
		uowRepo:     uowRepo,
	}
}

func (s *accountService) CreateAccount(ctx context.Context, req request.CreateAccountRequest) (*entities.AccountEntity, error) {

	result, err := s.uowRepo.Execute(ctx, func(ctx context.Context) (any, error) {
		account := &entities.AccountEntity{
			PublicAddress: req.PublicAddress,
			KycToken:      &req.KycToken,
		}

		if err := s.accountRepo.CreateAccount(ctx, account); err != nil {
			return &entities.AccountEntity{}, err
		}

		balance := &entities.BalanceEntity{
			AccountID:  account.ID,
			THBAmount:  0,
			USDTAmount: 0,
		}

		if err := s.balanceRepo.CreateBalance(ctx, balance); err != nil {
			return &entities.AccountEntity{}, err
		}

		return account, nil
	})

	if err != nil {
		msg := utils.FormatValidationError(err)
		return &entities.AccountEntity{}, entities.NewAppError(entities.ErrTypeInternal, msg, err)
	}

	return result.(*entities.AccountEntity), nil
}

func (s *accountService) GetAccounts(page int, limit int) ([]entities.AccountEntity, int64, error) {
	var accounts []entities.AccountEntity
	var total int64
	var errList, errCount error
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		accounts, errList = s.accountRepo.GetAccounts(page, limit)
	}()

	go func() {
		defer wg.Done()
		total, errCount = s.accountRepo.CountAccounts()
	}()

	wg.Wait()

	if errList != nil {
		msg := utils.FormatValidationError(errList)
		return []entities.AccountEntity{}, 0, entities.NewAppError(entities.ErrTypeInternal, msg, errList)
	}

	if errCount != nil {
		msg := utils.FormatValidationError(errCount)
		return []entities.AccountEntity{}, 0, entities.NewAppError(entities.ErrTypeInternal, msg, errCount)
	}

	return accounts, total, nil
}

func (s *accountService) GetAccountByID(id int) (*entities.AccountEntity, error) {
	account, err := s.accountRepo.GetAccountByID(id)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			return &entities.AccountEntity{}, entities.NewAppError(entities.ErrTypeNotFound, "account not found", err)
		}
		return &entities.AccountEntity{}, entities.NewAppError(entities.ErrTypeInternal, "internal server error", err)
	}
	return account, nil
}
