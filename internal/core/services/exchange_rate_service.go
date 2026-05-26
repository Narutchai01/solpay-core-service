package services

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
)

type ExchangeRateService interface {
	GetExchangeRate(symbol string) (*entities.ExchangeRate, error)
}

type exchangeRateService struct {
	repo ports.ExchangeRepository
}

func NewExchangeRateService(repo ports.ExchangeRepository) ExchangeRateService {
	return &exchangeRateService{
		repo: repo,
	}
}

func (s *exchangeRateService) GetExchangeRate(symbol string) (*entities.ExchangeRate, error) {
	rate, err := s.repo.GetExchangeRate(symbol)
	if err != nil {
		return nil, err
	}
	if len(*rate) == 0 {
		return nil, nil
	}
	return &(*rate)[0], nil
}
