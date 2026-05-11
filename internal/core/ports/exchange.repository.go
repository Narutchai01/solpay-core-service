package ports

import "github.com/Narutchai01/solpay-core-service/internal/entities"

type ExchangeRepository interface {
	GetExchangeRate(symbol string) (*[]entities.ExchangeRate, error)
}
