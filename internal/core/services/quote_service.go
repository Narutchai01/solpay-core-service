package services

import (
	"math"
	"time"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/dto/response"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
)

type QuoteService interface {
	CreateQuote(reqQuote request.CreateQuoteRequest, userID int64) (response.QuoteResponse, error)
	GetQuoteByID(id string) (*entities.Quote, error)
}

type quoteService struct {
	quoteRepo ports.QuoteRepository
}

func NewQuoteService(quoteRepo ports.QuoteRepository) QuoteService {
	return &quoteService{
		quoteRepo: quoteRepo,
	}
}

func (s *quoteService) CreateQuote(req request.CreateQuoteRequest, userID int64) (response.QuoteResponse, error) {

	var quote entities.Quote

	currentRate := 32.39

	rawUSDT := req.THBAmount / currentRate
	usdtAmount := math.Ceil(rawUSDT*1000000) / 1000000
	fee := math.Round(usdtAmount*0.01*1000000) / 1000000

	thbSatang := int64(math.Round(req.THBAmount * 100))
	expiresAt := time.Now().Add(60 * time.Second)

	quote = entities.Quote{
		AccountID:    userID, // สมมติว่าไม่มีการเชื่อมโยงกับบัญชีในตอนนี้
		Type:         req.ActionType,
		THBAmount:    thbSatang,
		USDTAmount:   usdtAmount,
		ExchangeRate: currentRate,
		ExpiresAt:    expiresAt,
		Status:       "PENDING",
		Fee:          fee,
	}

	err := s.quoteRepo.CreateQuote(&quote)
	if err != nil {
		return response.QuoteResponse{}, err
	}

	quoteResp := response.QuoteResponse{
		QuoteID:      quote.ID,
		THBAmount:    float64(quote.THBAmount) / 100, // แปลงกลับเป็นบาท
		USDTAmount:   quote.USDTAmount,
		ExchangeRate: quote.ExchangeRate,
		Fee:          quote.Fee,
	}

	return quoteResp, nil
}

func (s *quoteService) GetQuoteByID(id string) (*entities.Quote, error) {
	quote, err := s.quoteRepo.GetQuoteByID(id)
	if err != nil {
		return nil, err
	}
	return quote, nil
}
