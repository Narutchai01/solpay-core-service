package services

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/dto/response"
)

type SwapService interface {
	GetSwapQuote(query request.SwapQuoteRequest) (response.SwapQuoteData, error)
	BuildSwapUnsignedTransaction(req request.SwapUnsignedTransactionRequest, userID uint) (string, error)
}

type swapService struct {
	repo        ports.SwapRepository
	accountRepo ports.AccountRepository
}

func NewSwapService(repo ports.SwapRepository, accountRepo ports.AccountRepository) SwapService {
	return &swapService{
		repo:        repo,
		accountRepo: accountRepo,
	}
}

func (s *swapService) GetSwapQuote(query request.SwapQuoteRequest) (response.SwapQuoteData, error) {
	resp, err := s.repo.GetSwapQuote(query)
	if err != nil {
		return response.SwapQuoteData{}, err
	}
	return resp.Data, nil
}

func (s *swapService) BuildSwapUnsignedTransaction(req request.SwapUnsignedTransactionRequest, userID uint) (string, error) {
	account, err := s.accountRepo.GetAccountByID(int(userID))
	if err != nil {
		return "", err
	}

	resp, err := s.repo.BuildSwapUnsignedTransaction(req, account.PublicAddress)
	if err != nil {
		return "", err
	}
	return resp.Data.Transaction, nil
}
