package services

import (
	portRepo "github.com/Narutchai01/solpay-core-service/internal/core/ports/repositories"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/models/request"
)

type AccountService interface {
	CreateAccount(req request.CreateAccountRequest) (entities.AccountEntity, error)
}

type accountService struct {
	accountRepo portRepo.AccountRepository
}

func NewAccountService(accountRepo portRepo.AccountRepository) AccountService {
	return &accountService{
		accountRepo: accountRepo,
	}
}

func (s *accountService) CreateAccount(req request.CreateAccountRequest) (entities.AccountEntity, error) {
	account := entities.AccountEntity{
		PublicAddress: req.PublicAddress,
		KycToken:      req.KycToken,
	}

	createdAccount, err := s.accountRepo.CreateAccount(account)
	if err != nil {
		return entities.AccountEntity{}, err
	}
	return createdAccount, nil
}
