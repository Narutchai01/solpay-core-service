package ports

import (
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/dto/response"
)

type SwapRepository interface {
	GetSwapQuote(query request.SwapQuoteRequest) (response.SwapQuoteFullResponse, error)
}
