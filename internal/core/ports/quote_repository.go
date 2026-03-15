package ports

import "github.com/Narutchai01/solpay-core-service/internal/entities"

type QuoteRepository interface {
	CreateQuote(quote *entities.Quote) error
}
