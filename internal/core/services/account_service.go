package services

import (
	"errors"
	"fmt"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports/repositories"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/models/request"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
)

// Note: Define the AccountService interface
type AccountService interface {
	CreateAccount(req request.CreateAccountRequest) (entities.AccountEntity, error)
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

func (s *accountService) CreateAccount(req request.CreateAccountRequest) (entities.AccountEntity, error) {
	// NOTE: You might want to add additional validation or business logic here
	account := entities.AccountEntity{
		PublicAddress: req.PublicAddress,
		KycToken:      &req.KycToken,
	}

	// NOTE: Handle potential conflicts or errors during account creation
	createdAccount, err := s.accountRepo.CreateAccount(account)
	if err != nil {
		// Note: Handle specific error cases
		if errors.Is(err, entities.ErrConflict) {
			// NOTE : if account already exists
			errMessage := fmt.Sprintf("Account with public address %s already exists", req.PublicAddress)
			return entities.AccountEntity{}, entities.NewAppError(entities.ErrTypeConflict, errMessage, err)
		}
		msg := utils.FormatValidationError(err)
		// Generic error handling
		return entities.AccountEntity{}, entities.NewAppError(entities.ErrTypeInternal, msg, err)
	}
	return createdAccount, nil
}
