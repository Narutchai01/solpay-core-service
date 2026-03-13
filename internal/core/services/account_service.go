package services

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
)

// AccountService defines operations for managing accounts.
type AccountService interface {
	CreateAccount(ctx context.Context, req request.CreateAccountRequest) (*entities.AccountEntity, error)
	GetAccounts(page, limit int) ([]entities.AccountEntity, int64, error)
	GetAccountByID(id int) (*entities.AccountEntity, error)
}

type accountService struct {
	accountRepo ports.AccountRepository
	balanceRepo ports.BalanceRepository
	uow         ports.UnitOfWork
}

// NewAccountService creates a new AccountService.
func NewAccountService(accountRepo ports.AccountRepository, balanceRepo ports.BalanceRepository, uow ports.UnitOfWork) AccountService {
	return &accountService{
		accountRepo: accountRepo,
		balanceRepo: balanceRepo,
		uow:         uow,
	}
}

func (s *accountService) CreateAccount(ctx context.Context, req request.CreateAccountRequest) (*entities.AccountEntity, error) {
	// Return existing account if the address is already registered.
	existing, err := s.accountRepo.GetAccountByPublicAddress(req.PublicAddress)
	if err == nil && existing.ID != 0 {
		return existing, nil
	}
	if err != nil && !errors.Is(err, entities.ErrNotFound) {
		return nil, entities.NewAppError(entities.ErrTypeInternal, "failed to check existing account", err)
	}

	result, err := s.uow.Execute(ctx, func(txCtx context.Context) (any, error) {
		account := &entities.AccountEntity{
			PublicAddress: req.PublicAddress,
		}
		if err := s.accountRepo.CreateAccount(txCtx, account); err != nil {
			return nil, err
		}

		balance := &entities.BalanceEntity{
			AccountID:  account.ID,
			THBAmount:  0,
			USDTAmount: 0,
		}
		if err := s.balanceRepo.CreateBalance(txCtx, balance); err != nil {
			return nil, err
		}

		return account, nil
	})
	if err != nil {
		return nil, entities.NewAppError(entities.ErrTypeInternal, "failed to create account", err)
	}

	return result.(*entities.AccountEntity), nil
}

func (s *accountService) GetAccounts(page, limit int) ([]entities.AccountEntity, int64, error) {
	var (
		accounts []entities.AccountEntity
		total    int64
		errList  error
		errCount error
		wg       sync.WaitGroup
	)

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

	if err := errors.Join(errList, errCount); err != nil {
		return nil, 0, entities.NewAppError(entities.ErrTypeInternal, "failed to list accounts", err)
	}

	return accounts, total, nil
}

func (s *accountService) GetAccountByID(id int) (*entities.AccountEntity, error) {
	account, err := s.accountRepo.GetAccountByID(id)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			return nil, entities.NewAppError(entities.ErrTypeNotFound, fmt.Sprintf("account %d not found", id), err)
		}
		return nil, entities.NewAppError(entities.ErrTypeInternal, "failed to get account", err)
	}
	return account, nil
}
