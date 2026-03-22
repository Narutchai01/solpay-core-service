package services

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/dto/response"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/models"
)

func isInvalidSolanaAddressError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, "invalid sender address") ||
		strings.Contains(errMsg, "invalid receive address") ||
		strings.Contains(errMsg, "invalid mint token address")
}

type QuoteService interface {
	CreateQuote(reqQuote request.CreateQuoteRequest, accountID uint) (response.QuoteResponse, error)
	GetQuoteByID(id string) (*entities.Quote, error)
	ConFirmQuote(id string, accountID uint) (string, error)
}

type quoteService struct {
	quoteRepo   ports.QuoteRepository
	solanaRepo  ports.SolanaClient
	accountRepo ports.AccountRepository
}

func NewQuoteService(quoteRepo ports.QuoteRepository, solanaRepo ports.SolanaClient, accountRepo ports.AccountRepository) QuoteService {
	return &quoteService{
		quoteRepo:   quoteRepo,
		solanaRepo:  solanaRepo,
		accountRepo: accountRepo,
	}
}

func (s *quoteService) CreateQuote(req request.CreateQuoteRequest, accountID uint) (response.QuoteResponse, error) {

	if req.ActionType == string(entities.ONCHAIN) {
		if req.PromptPayID == "" {
			return response.QuoteResponse{}, entities.NewAppError(entities.ErrTypeBadRequest, "promptpay_id is required for onchain transactions", nil)
		}
	}

	var quote entities.Quote

	currentRate := 32.39

	rawUSDT := req.THBAmount / currentRate
	usdtAmount := math.Ceil(rawUSDT*1000000) / 1000000
	fee := math.Round(usdtAmount*0.01*1000000) / 1000000

	thbSatang := int64(math.Round(req.THBAmount * 100))
	expiresAt := time.Now().Add(60 * time.Second)

	quote = entities.Quote{
		AccountID:    accountID, // สมมติว่าไม่มีการเชื่อมโยงกับบัญชีในตอนนี้
		Type:         req.ActionType,
		THBAmount:    thbSatang,
		USDTAmount:   usdtAmount,
		ExchangeRate: currentRate,
		ExpiresAt:    expiresAt,
		PromptPayID:  &req.PromptPayID,
		Status:       string(entities.ACTIVE),
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

func (s *quoteService) ConFirmQuote(id string, accountID uint) (string, error) {
	if s.solanaRepo == nil {
		return "", entities.NewAppError(entities.ErrTypeInternal, "solana client is not configured", nil)
	}

	quote, err := s.quoteRepo.GetQuoteByID(id)
	if err != nil {
		return "", err
	}

	if time.Now().After(quote.ExpiresAt) {
		return "", entities.NewAppError(entities.ErrTypeBadRequest, "quote has expired", nil)
	}

	account, err := s.accountRepo.GetAccountByID(int(accountID))
	if err != nil {
		return "", err
	}

	if quote.AccountID != account.ID {
		return "", entities.NewAppError(entities.ErrTypeConflict, "quote does not belong to the account", nil)
	}

	rawTx := models.BuildTXUnsigned{
		SenderAddress: account.PublicAddress,
		Amount:        uint64(quote.USDTAmount * 1e6), // แปลงเป็นหน่วยเล็กสุดของ USDT
		Decimals:      6,                              // สมมติว่า USDT มี 6 ทศนิยม
	}

	txHashBase64, err := s.solanaRepo.BuildUnsignedTransfer(context.Background(), rawTx)
	if err != nil {
		if isInvalidSolanaAddressError(err) {
			return "", entities.NewAppError(entities.ErrTypeBadRequest, "invalid Solana address configuration", err)
		}

		var appErr *entities.AppError
		if errors.As(err, &appErr) {
			return "", appErr
		}

		return "", entities.NewAppError(entities.ErrTypeInternal, "failed to build unsigned transfer", err)
	}

	if txHashBase64 == "" {
		return "", entities.NewAppError(entities.ErrTypeInternal, "empty unsigned transaction payload", nil)
	}

	quote.Status = string(entities.USED)

	err = s.quoteRepo.UpdateQuote(quote)
	if err != nil {
		return "", err
	}

	return txHashBase64, nil
}
