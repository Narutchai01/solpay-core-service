package services

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/dto/response"
)

type SwapService interface {
	GetSwapQuote(query request.SwapQuoteRequest) (response.SwapQuoteData, error)
}

type swapService struct {
	repo ports.SwapRepository
}

func NewSwapService(repo ports.SwapRepository) SwapService {
	return &swapService{repo: repo}
}

func (s *swapService) GetSwapQuote(query request.SwapQuoteRequest) (response.SwapQuoteData, error) {
	resp, err := s.repo.GetSwapQuote(query)
	if err != nil {
		return response.SwapQuoteData{}, err
	}
	return resp.Data, nil
}
