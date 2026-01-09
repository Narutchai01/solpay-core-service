package services

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports/repositories"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
)

// Note: Define the AccountService interface
type AccountService interface {
	CreateAccount(req request.CreateAccountRequest) (*entities.AccountEntity, error)
	GetAccounts(page int, limit int) ([]entities.AccountEntity, int64, error)
	GetAccountByID(id int) (*entities.AccountEntity, error)
	GetAccountByPublicAddress(address string) (*entities.AccountEntity, error)
}

// Note: Implement the AccountService interface
type accountService struct {
	accountRepo repositories.AccountRepository
}

// Note: Constructor function for AccountService
func NewAccountService(accountRepo repositories.AccountRepository) AccountService {
	return &accountService{
		accountRepo: accountRepo,
	}
}

func (s *accountService) CreateAccount(req request.CreateAccountRequest) (*entities.AccountEntity, error) {
	// NOTE: You might want to add additional validation or business logic here
	account := entities.AccountEntity{
		PublicAddress: req.PublicAddress,
		KycToken:      &req.KycToken,
	}

	// NOTE: Handle potential conflicts or errors during account creation
	err := s.accountRepo.CreateAccount(&account)
	if err != nil {
		// Note: Handle specific error cases
		if errors.Is(err, entities.ErrConflict) {
			// NOTE : if account already exists
			errMessage := fmt.Sprintf("Account with public address %s already exists", req.PublicAddress)
			return nil, entities.NewAppError(entities.ErrTypeConflict, errMessage, err)
		}
		msg := utils.FormatValidationError(err)
		// Generic error handling
		return nil, entities.NewAppError(entities.ErrTypeInternal, msg, err)
	}
	return &account, nil
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
	var account entities.AccountEntity
	account, err := s.accountRepo.GetAccountByID(id)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			return &entities.AccountEntity{}, entities.NewAppError(entities.ErrTypeNotFound, "account not found", err)
		}
		return &entities.AccountEntity{}, entities.NewAppError(entities.ErrTypeInternal, "internal server error", err)
	}
	return &account, nil
}

func (s *accountService) GetAccountByPublicAddress(address string) (*entities.AccountEntity, error) {
	var account entities.AccountEntity
	account, err := s.accountRepo.GetAccountByPublicAddress(address)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			return &entities.AccountEntity{}, entities.NewAppError(entities.ErrTypeNotFound, "account not found", err)
		}
		return &entities.AccountEntity{}, entities.NewAppError(entities.ErrTypeInternal, "internal server error", err)
	}
	return &account, nil
}
