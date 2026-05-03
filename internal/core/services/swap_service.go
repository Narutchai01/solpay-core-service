package services

import (
	"context"
	"errors"
	"math"
	"strconv"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/dto/response"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/google/uuid"
)

type SwapService interface {
	GetSwapQuote(query request.SwapQuoteRequest) (response.SwapQuoteData, error)
	BuildSwapUnsignedTransaction(req request.SwapUnsignedTransactionRequest, userID uint) (response.BuildSwapTransactionResponse, error)
	ExecuteSwapTransaction(ctx context.Context, req request.ExecuteSwapTransactionRequest, accountID uint) (response.TransactionDTO, error)
}

type swapService struct {
	repo            ports.SwapRepository
	accountRepo     ports.AccountRepository
	transactionRepo ports.TransactionRepository
	uow             ports.UnitOfWork
	publisher       ports.Publisher
}

func NewSwapService(
	repo ports.SwapRepository,
	accountRepo ports.AccountRepository,
	transactionRepo ports.TransactionRepository,
	uow ports.UnitOfWork,
	publisher ports.Publisher,
) SwapService {
	return &swapService{
		repo:            repo,
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		uow:             uow,
		publisher:       publisher,
	}
}

func (s *swapService) GetSwapQuote(query request.SwapQuoteRequest) (response.SwapQuoteData, error) {
	// Convert AmountIn from decimal string to raw amount string (assuming 9 decimals)
	if query.AmountIn != "" {
		amountDecimal, err := strconv.ParseFloat(query.AmountIn, 64)
		if err == nil {
			rawAmount := int64(math.Round(amountDecimal * 1e9))
			query.AmountIn = strconv.FormatInt(rawAmount, 10)
		}
	}

	resp, err := s.repo.GetSwapQuote(query)
	if err != nil {
		return response.SwapQuoteData{}, entities.NewAppError(entities.ErrTypeInternal, "failed to get swap quote", err)
	}
	return resp.Data, nil
}

func (s *swapService) BuildSwapUnsignedTransaction(req request.SwapUnsignedTransactionRequest, userID uint) (response.BuildSwapTransactionResponse, error) {
	// Convert AmountIn from decimal string to raw amount string (assuming 9 decimals)
	if req.AmountIn != "" {
		amountDecimal, err := strconv.ParseFloat(req.AmountIn, 64)
		if err == nil {
			rawAmount := int64(math.Round(amountDecimal * 1e9))
			req.AmountIn = strconv.FormatInt(rawAmount, 10)
		}
	}

	account, err := s.accountRepo.GetAccountByID(int(userID))
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			return response.BuildSwapTransactionResponse{}, entities.NewAppError(entities.ErrTypeNotFound, "account not found", err)
		}
		return response.BuildSwapTransactionResponse{}, entities.NewAppError(entities.ErrTypeInternal, "failed to get account", err)
	}

	resp, err := s.repo.BuildSwapUnsignedTransaction(req, account.PublicAddress)
	if err != nil {
		return response.BuildSwapTransactionResponse{}, entities.NewAppError(entities.ErrTypeInternal, "failed to build swap transaction", err)
	}
	return response.BuildSwapTransactionResponse{
		Transaction: resp.Data.Transaction,
	}, nil
}

func (s *swapService) ExecuteSwapTransaction(ctx context.Context, req request.ExecuteSwapTransactionRequest, accountID uint) (response.TransactionDTO, error) {
	cfg := config.LoadConfig()
	txUUID, err := uuid.NewV7()
	if err != nil {
		return response.TransactionDTO{}, entities.NewAppError(entities.ErrTypeInternal, "failed to generate transaction UUID", err)
	}

	tx := &entities.TransactionEntity{
		TransactionUUID: txUUID,
		AccountID:       accountID,
		TransactionType: string(entities.SWAP),
		USDTAmount:      func() float64 { v, _ := strconv.ParseFloat(req.USDTAmount, 64); return v }(),
		SOLAmount:       func() float64 { v, _ := strconv.ParseFloat(req.SOLAmount, 64); return v }(),
	}

	result, err := s.uow.Execute(ctx, func(txCtx context.Context) (any, error) {
		if err := s.transactionRepo.CreateTransaction(txCtx, tx); err != nil {
			return nil, err
		}

		signature := signatureFromTxHash(req.TxHash)

		txOnChain := &entities.TransactionOnChain{
			TransactionID: tx.TransactionUUID,
			Signature:     signature,
			TxHash:        req.TxHash,
		}

		if err := s.transactionRepo.CreateTransactionOnChain(txCtx, txOnChain); err != nil {
			return nil, err
		}
		return tx, nil
	})

	if err != nil {
		return response.TransactionDTO{}, entities.NewAppError(entities.ErrTypeInternal, "failed to execute swap transaction", err)
	}

	tx = result.(*entities.TransactionEntity)

	s.publisher.PublishTransactionMessage(cfg.TRANSACTION_ORCHESTRATOR_QUEUE, tx, string(entities.StatusSolanaSubmitted))

	return *response.FormatTransactionDTO(tx), nil
}
