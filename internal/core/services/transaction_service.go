package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/dto/response"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/infra/supabase"
	"github.com/Narutchai01/solpay-core-service/internal/models"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"github.com/google/uuid"
)

const errPublisherNotConfigured = "publisher is not configured"
const errUpdateTransactionStatusFmt = "update transaction status: %w"

func satangToTHB(amount float64) float64 {
	return math.Round((amount/100.0)*100) / 100
}

// TransactionService defines operations for managing transactions.
type TransactionService interface {
	CreateTransaction(ctx context.Context, req request.CreateTransactionRequest) (*entities.TransactionEntity, error)
	GetTransactionByUUID(txUUID uuid.UUID) (*entities.TransactionEntity, error)
	QueryTransactionSummary(ctx context.Context, month, year int) (*response.TransactionChartSummary, error)
	HandleTransactionUpdate(ctx context.Context, msg []byte) error
	GetTransactions(query request.TransactionQuery, accountID *uint) ([]entities.TransactionEntity, int64, error)
	GetSpendingSummary(ctx context.Context, accountID uint) (*response.OverallSpendingSummaryDTO, error)
	GetMonthlySpendingSummary(ctx context.Context, accountID uint) ([]response.MonthlySpendingDTO, error)
}

type transactionService struct {
	transactionRepo ports.TransactionRepository
	uow             ports.UnitOfWork
	publisher       ports.Publisher
	wsHub           ports.WebSocketPort
	cfg             *config.Config
}

// NewTransactionService creates a new TransactionService.
func NewTransactionService(
	transactionRepo ports.TransactionRepository,
	uow ports.UnitOfWork,
	publisher ports.Publisher,
	wsHub ports.WebSocketPort,
	cfg *config.Config,
) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		uow:             uow,
		publisher:       publisher,
		wsHub:           wsHub,
		cfg:             cfg,
	}
}

// ---------------------------------------------------------------------------
// CreateTransaction
// ---------------------------------------------------------------------------

func (s *transactionService) CreateTransaction(ctx context.Context, req request.CreateTransactionRequest) (*entities.TransactionEntity, error) {
	txUUID, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generate UUID: %w", err)
	}

	transaction := &entities.TransactionEntity{
		TransactionUUID: txUUID,
		AccountID:       1,
		TransactionType: req.TransactionType,
		THBAmount:       req.THBAmount,
		USDTAmount:      req.USDTAmount,
		Fee:             req.Fee,
	}

	result, err := s.uow.Execute(ctx, func(txCtx context.Context) (any, error) {
		if err := s.transactionRepo.CreateTransaction(txCtx, transaction); err != nil {
			return nil, err
		}
		if err := s.createSubTransactions(txCtx, req, txUUID); err != nil {
			return nil, err
		}
		return transaction, nil
	})
	if err != nil {
		return nil, fmt.Errorf("create transaction: %w", err)
	}

	tx := result.(*entities.TransactionEntity)

	return tx, nil
}

// createSubTransactions creates the type-specific child records within the UoW.
func (s *transactionService) createSubTransactions(ctx context.Context, req request.CreateTransactionRequest, txID uuid.UUID) error {
	switch req.TransactionType {
	case string(entities.TOPUP):
		return s.createTransactionOnChain(ctx, req, txID)
	case string(entities.OFFCHAIN):
		return s.createTransactionOffChain(ctx, req, txID)
	case string(entities.ONCHAIN):
		if err := s.createTransactionOnChain(ctx, req, txID); err != nil {
			return err
		}
		return s.createTransactionOffChain(ctx, req, txID)
	default:
		return fmt.Errorf("unsupported transaction type: %s", req.TransactionType)
	}
}

func (s *transactionService) createTransactionOnChain(ctx context.Context, req request.CreateTransactionRequest, txID uuid.UUID) error {
	if req.TxHash == nil {
		return errors.New("tx_hash is required for on-chain transactions")
	}

	signature := signatureFromTxHash(*req.TxHash)

	return s.transactionRepo.CreateTransactionOnChain(ctx, &entities.TransactionOnChain{
		TransactionID: txID,
		Signature:     signature,
		TxHash:        *req.TxHash,
	})
}

func (s *transactionService) createTransactionOffChain(ctx context.Context, req request.CreateTransactionRequest, txID uuid.UUID) error {
	if req.PromptPayID == nil {
		return errors.New("prompt_pay_id is required for off-chain transactions")
	}
	return s.transactionRepo.CreateTransactionOffChain(ctx, &entities.TransactionOffChain{
		TransactionID: txID,
		PromptPayID:   *req.PromptPayID,
	})
}

// ---------------------------------------------------------------------------
// GetTransactionByID
// ---------------------------------------------------------------------------

func (s *transactionService) GetTransactionByUUID(txUUID uuid.UUID) (*entities.TransactionEntity, error) {
	tx, err := s.transactionRepo.GetTransactionByUUID(txUUID)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			return nil, entities.NewAppError(entities.ErrTypeNotFound, fmt.Sprintf("transaction %s not found", txUUID), err)
		}
		return nil, entities.NewAppError(entities.ErrTypeInternal, "failed to get transaction", err)
	}
	return tx, nil
}

func (s *transactionService) HandleTransactionUpdate(ctx context.Context, msg []byte) error {
	var event request.TransactionMessage
	if err := json.Unmarshal(msg, &event); err != nil {
		return fmt.Errorf("unmarshal transaction message: %w", err)
	}

	if event.TxID == "" || event.Status == "" {
		return entities.NewAppError(entities.ErrTypeBadRequest, "tx_id and status are required", nil)
	}

	txUUID, err := uuid.Parse(event.TxID)
	if err != nil {
		return fmt.Errorf("parse transaction UUID: %w", err)
	}

	tx, err := s.transactionRepo.GetTransactionByUUID(txUUID)
	if err != nil {
		return fmt.Errorf("get transaction by UUID: %w", err)
	}

	if tx.Status == event.Status {
		return nil
	}

	if err := s.updateStatusAndNotify(ctx, event.TxID, event.Status); err != nil {
		return err
	}

	if err := s.routeEvent(ctx, tx, event); err != nil {
		return s.finalizeFailed(ctx, event.TxID, err)
	}

	return nil
}

func (s *transactionService) finalizeFailed(ctx context.Context, txID string, err error) error {
	log.Printf("finalizing transaction %s as FAILED due to: %v", txID, err)
	if updateErr := s.updateStatusAndNotify(ctx, txID, string(entities.StatusFailed)); updateErr != nil {
		return fmt.Errorf("failed to finalize transaction status to FAILED: %v (original error: %w)", updateErr, err)
	}
	return err
}

func (s *transactionService) routeEvent(ctx context.Context, tx *entities.TransactionEntity, event request.TransactionMessage) error {
	switch tx.TransactionType {
	case string(entities.TOPUP):
		return s.topupWorkflow(ctx, tx, event)
	case string(entities.OFFCHAIN):
		return s.offchainWorkflow(ctx, tx, event)
	case string(entities.ONCHAIN):
		return s.onchainWorkflow(ctx, tx, event)
	case string(entities.SWAP):
		return s.swapWorkflow(ctx, tx, event)
	default:
		return nil
	}
}

func (s *transactionService) topupWorkflow(ctx context.Context, tx *entities.TransactionEntity, event request.TransactionMessage) error {
	switch event.Status {
	case string(entities.StatusSolanaSubmitted):
		return s.publishBlockchainTransaction(tx)
	case string(entities.StatusSolanaFailed), string(entities.StatusBalanceFailed), string(entities.StatusFailed):
		return s.updateStatusAndNotify(ctx, tx.TransactionUUID.String(), string(entities.StatusFailed))
	case string(entities.StatusSolanaSuccess):
		return s.publishBalanceTransaction(tx, string(entities.ActionDeposit))
	case string(entities.StatusBalanceUpdated):
		return s.updateStatusAndNotify(ctx, tx.TransactionUUID.String(), string(entities.StatusCompleted))
	default:
		return fmt.Errorf("unhandled transaction status in topup workflow: %s", event.Status)
	}
}

func (s *transactionService) offchainWorkflow(ctx context.Context, tx *entities.TransactionEntity, event request.TransactionMessage) error {
	switch event.Status {
	case string(entities.StatusBalanceWithdrawing):
		return s.publishBalanceTransaction(tx, string(entities.ActionWithdraw))
	case string(entities.StatusBalanceUpdating):
		return s.publishBalanceTransaction(tx, string(entities.ActionWithdraw))
	case string(entities.StatusBalanceUpdated):
		if err := s.updateStatusAndNotify(ctx, tx.TransactionUUID.String(), string(entities.StatusBalanceUpdated)); err != nil {
			return err
		}
		return s.publishPaymentTransaction(tx)
	case string(entities.StatusPaymentSuccess):
		cfg := s.cfg

		// Generate Slip
		var address string
		if tx.Account != nil {
			address = tx.Account.PublicAddress
		}
		var promptPayID string
		if tx.TransactionOffChain != nil {
			promptPayID = tx.TransactionOffChain.PromptPayID
		}

		slipData := utils.SlipOffchain{
			Address:       address,
			Amount:        tx.THBAmount / 100.0, // Amount is stored as satang, but GetSlipOFFCHAINInformation seems to expect float64 THB
			TransactionID: tx.TransactionUUID.String(),
			PromptPayID:   promptPayID,
			CreatedAt:     tx.CreatedAt.Format("02/01/2006 15:04:05"),
		}

		slipBytes, err := utils.GetSlipOFFCHAINInformation(slipData)
		if err != nil {
			return fmt.Errorf("generate slip: %w", err)
		}

		// Upload to Supabase
		supabaseStorage := supabase.NewSupabaseStorage(cfg.SUPABASE_PRIVATE_KEY, cfg.SUPABASE_URL)
		fileName := fmt.Sprintf("%s.png", tx.TransactionUUID.String())
		slipURL, err := supabaseStorage.UploadFile("slip", fileName, slipBytes)
		if err != nil {
			return fmt.Errorf("upload slip: %w", err)
		}

		// Update TransactionOffChain with SlipURL
		if tx.TransactionOffChain != nil {
			tx.TransactionOffChain.SlipURL = &slipURL
			if err := s.transactionRepo.UpdateTransactionOffChain(ctx, tx.TransactionOffChain); err != nil {
				return fmt.Errorf("update transaction offchain: %w", err)
			}
		}

		return s.updateStatusAndNotify(ctx, tx.TransactionUUID.String(), string(entities.StatusCompleted))
	case string(entities.StatusBalanceFailed), string(entities.StatusPaymentFailed), string(entities.StatusFailed):
		return s.updateStatusAndNotify(ctx, tx.TransactionUUID.String(), string(entities.StatusFailed))
	default:
		return fmt.Errorf("unhandled transaction status in offchain workflow: %s", event.Status)
	}
}

func (s *transactionService) onchainWorkflow(ctx context.Context, tx *entities.TransactionEntity, event request.TransactionMessage) error {
	switch event.Status {
	case string(entities.StatusSolanaSubmitted):
		return s.publishBlockchainTransaction(tx)
	case string(entities.StatusSolanaFailed), string(entities.StatusBalanceFailed), string(entities.StatusPaymentFailed), string(entities.StatusFailed):
		return s.updateStatusAndNotify(ctx, tx.TransactionUUID.String(), string(entities.StatusFailed))
	case string(entities.StatusSolanaSuccess):
		return s.publishPaymentTransaction(tx)
	case string(entities.StatusPaymentSuccess):
		cfg := s.cfg

		// Generate Slip
		var address string
		if tx.Account != nil {
			address = tx.Account.PublicAddress
		}
		var promptPayID string
		if tx.TransactionOffChain != nil {
			promptPayID = tx.TransactionOffChain.PromptPayID
		}

		slipData := utils.SlipOnchain{
			Address:       address,
			THBAmount:     tx.THBAmount / 100.0, // convert satang to THB
			USDTAmount:    tx.USDTAmount,
			FreeAmount:    tx.Fee,
			TransactionID: tx.TransactionUUID.String(),
			PromptPayID:   promptPayID,
			CreatedAt:     tx.CreatedAt.Format("02/01/2006 15:04:05"),
		}

		slipBytes, err := utils.GetSlipOnChain(slipData)
		if err != nil {
			return fmt.Errorf("generate onchain slip: %w", err)
		}

		// Upload to Supabase
		supabaseStorage := supabase.NewSupabaseStorage(cfg.SUPABASE_PRIVATE_KEY, cfg.SUPABASE_URL)
		fileName := fmt.Sprintf("onchain_%s.png", tx.TransactionUUID.String())
		slipURL, err := supabaseStorage.UploadFile("slip", fileName, slipBytes)
		if err != nil {
			return fmt.Errorf("upload onchain slip: %w", err)
		}

		// Update TransactionOffChain with SlipURL
		if tx.TransactionOffChain != nil {
			tx.TransactionOffChain.SlipURL = &slipURL
			if err := s.transactionRepo.UpdateTransactionOffChain(ctx, tx.TransactionOffChain); err != nil {
				return fmt.Errorf("update transaction offchain: %w", err)
			}
		}

		return s.updateStatusAndNotify(ctx, tx.TransactionUUID.String(), string(entities.StatusCompleted))
	default:
		return fmt.Errorf("unhandled transaction status in onchain workflow: %s", event.Status)
	}
}

func (s *transactionService) swapWorkflow(ctx context.Context, tx *entities.TransactionEntity, event request.TransactionMessage) error {
	switch event.Status {
	case string(entities.StatusSolanaSubmitted):
		return s.publishBlockchainTransaction(tx)
	case string(entities.StatusSolanaSuccess):
		// For SWAP, success from blockchain means the transaction is complete
		return s.updateStatusAndNotify(ctx, tx.TransactionUUID.String(), string(entities.StatusCompleted))
	case string(entities.StatusSolanaFailed), string(entities.StatusBalanceFailed), string(entities.StatusPaymentFailed), string(entities.StatusFailed):
		return s.updateStatusAndNotify(ctx, tx.TransactionUUID.String(), string(entities.StatusFailed))
	default:
		return fmt.Errorf("unhandled transaction status in swap workflow: %s", event.Status)
	}
}

func (s *transactionService) updateStatusAndNotify(ctx context.Context, txID string, status string) error {
	if err := s.transactionRepo.UpdateTransactionStatus(ctx, txID, status); err != nil {
		return fmt.Errorf(errUpdateTransactionStatusFmt, err)
	}

	if s.wsHub == nil {
		return nil
	}

	if status != string(entities.StatusFailed) && status != string(entities.StatusCompleted) {
		return nil
	}

	txUUID, err := uuid.Parse(txID)
	if err != nil {
		log.Printf("failed to parse tx id for websocket notify: %v", err)
		return nil
	}

	tx, err := s.transactionRepo.GetTransactionByUUID(txUUID)
	if err != nil {
		log.Printf("failed to fetch transaction for websocket notify: %v", err)
		return nil
	}

	jsonPayload, err := json.Marshal(response.FormatTransactionDTO(tx))
	if err != nil {
		log.Printf("failed to marshal websocket transaction payload: %v", err)
		return nil
	}

	s.wsHub.NotifyTransactionStatus(txID, jsonPayload)
	return nil
}

func (s *transactionService) publishBlockchainTransaction(tx *entities.TransactionEntity) error {
	if s.publisher == nil {
		return entities.NewAppError(entities.ErrTypeInternal, errPublisherNotConfigured, nil)
	}
	if tx.TransactionOnChain == nil || tx.TransactionOnChain.TxHash == "" {
		return entities.NewAppError(entities.ErrTypeBadRequest, "missing on-chain transaction payload", nil)
	}

	cfg := s.cfg

	msg := models.SolanaTxMessage{
		TxID:     tx.TransactionUUID.String(),
		Base64Tx: tx.TransactionOnChain.TxHash,
		MetaData: models.SolanaMetaData{
			TransactionType: tx.TransactionType,
			AmountTHB:       int(tx.THBAmount),
			AmountUSDC:      int(tx.USDTAmount),
			AccountID:       int(tx.AccountID),
		},
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal solana tx message: %w", err)
	}

	return s.publisher.Publish(cfg.SOLANA_WORK_QUEUE, jsonMsg)
}

func (s *transactionService) publishBalanceTransaction(tx *entities.TransactionEntity, action string) error {
	if s.publisher == nil {
		return entities.NewAppError(entities.ErrTypeInternal, errPublisherNotConfigured, nil)
	}

	cfg := s.cfg
	msg := request.UpdateBalanceCommand{
		AccountID:     uint(tx.AccountID),
		TransactionID: tx.TransactionUUID.String(),
		THBAmount:     int64(tx.THBAmount),
		Action:        action,
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal balance message: %w", err)
	}

	err = s.publisher.Publish(cfg.BALANCE_QUEUE, jsonMsg)
	if err != nil {
		return fmt.Errorf("publish balance message: %w", err)
	}

	return nil
}

func (s *transactionService) publishPaymentTransaction(tx *entities.TransactionEntity) error {
	if s.publisher == nil {
		return entities.NewAppError(entities.ErrTypeInternal, errPublisherNotConfigured, nil)
	}
	if tx.TransactionOffChain == nil || tx.TransactionOffChain.PromptPayID == "" {
		return entities.NewAppError(entities.ErrTypeBadRequest, "missing off-chain transaction payload", nil)
	}

	cfg := s.cfg
	msg := request.RequestPaymentQueue{
		TransactionID: string(tx.TransactionUUID.String()),
		Amount:        int64(tx.THBAmount),
		Number:        tx.TransactionOffChain.PromptPayID,
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal balance message: %w", err)
	}

	err = s.publisher.Publish(cfg.PAYMENT_QUEUE, jsonMsg)
	if err != nil {
		return fmt.Errorf("publish balance message: %w", err)
	}

	return nil
}

func (s *transactionService) GetTransactions(query request.TransactionQuery, accountID *uint) ([]entities.TransactionEntity, int64, error) {
	total, err := s.transactionRepo.CountTransactions(query, accountID)
	if err != nil {
		return nil, 0, entities.NewAppError(entities.ErrTypeInternal, "failed to count transactions", err)
	}

	if total == 0 {
		return []entities.TransactionEntity{}, 0, nil
	}

	transactions, err := s.transactionRepo.GetTransactions(query, accountID)
	if err != nil {
		return nil, 0, entities.NewAppError(entities.ErrTypeInternal, "failed to get transactions", err)
	}

	return transactions, total, nil
}

func (s *transactionService) QueryTransactionSummary(ctx context.Context, month, year int) (*response.TransactionChartSummary, error) {
	rows, err := s.transactionRepo.QueryTransactionSummary(ctx, month, year)
	if err != nil {
		return nil, entities.NewAppError(entities.ErrTypeInternal, "failed to query transaction summary", err)
	}

	// Index raw rows by date+type for O(1) lookup
	type key struct{ date, txType string }
	indexedTHBAmount := make(map[key]float64, len(rows))
	indexedUSDTAmount := make(map[key]float64, len(rows))
	indexedFee := make(map[key]float64, len(rows))
	totalCompletedCount := 0
	for _, r := range rows {
		k := key{r.Date, r.TransactionType}
		indexedTHBAmount[k] = r.TotalTHBAmount
		indexedUSDTAmount[k] = r.TotalUSDTAmount
		indexedFee[k] = r.TotalFee
		totalCompletedCount += r.TotalCount
	}

	// Build a full calendar for the requested month (all days, no gaps)
	daysInMonth := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()

	chartData := make([]response.TransactionChartData, 0, daysInMonth)
	var totalDeposit, totalWithdraw, totalFee float64

	for day := 1; day <= daysInMonth; day++ {
		date := fmt.Sprintf("%04d-%02d-%02d", year, month, day)
		label := fmt.Sprintf("%02d", day)

		deposit := indexedUSDTAmount[key{date, string(entities.TOPUP)}] + indexedUSDTAmount[key{date, string(entities.ONCHAIN)}]
		withdraw := indexedTHBAmount[key{date, string(entities.OFFCHAIN)}] + indexedTHBAmount[key{date, string(entities.ONCHAIN)}]
		withdrawTHB := satangToTHB(withdraw)
		fee := indexedFee[key{date, string(entities.TOPUP)}] + indexedFee[key{date, string(entities.OFFCHAIN)}] + indexedFee[key{date, string(entities.ONCHAIN)}]

		totalDeposit += deposit
		totalWithdraw += withdraw
		totalFee += fee

		chartData = append(chartData, response.TransactionChartData{
			Date:     date,
			Label:    label,
			Deposit:  deposit,
			Withdraw: response.NewDecimal2(withdrawTHB),
			Fee:      fee,
		})
	}

	return &response.TransactionChartSummary{
		TotalDeposit:        totalDeposit,
		TotalWithdraw:       response.NewDecimal2(satangToTHB(totalWithdraw)),
		TotalFee:            totalFee,
		TotalCompletedCount: totalCompletedCount,
		ChartData:           chartData,
	}, nil
}

func (s *transactionService) GetSpendingSummary(ctx context.Context, accountID uint) (*response.OverallSpendingSummaryDTO, error) {
	now := time.Now()
	month := int(now.Month())
	year := now.Year()

	categorySummaries, err := s.transactionRepo.GetSpendingSummary(ctx, accountID, month, year)
	if err != nil {
		return nil, entities.NewAppError(entities.ErrTypeInternal, "failed to get spending summary", err)
	}

	var totalSpentCurrentMonth float64
	for _, cs := range categorySummaries {
		totalSpentCurrentMonth += cs.TotalSpent
	}

	categoryDtos := make([]response.SpendingSummaryDTO, 0, len(categorySummaries))
	for _, cs := range categorySummaries {
		name := cs.CategoryName
		if name == "" {
			name = "Uncategorized"
		}

		var percentage float64
		if totalSpentCurrentMonth > 0 {
			percentage = (cs.TotalSpent / totalSpentCurrentMonth) * 100
		}

		categoryDtos = append(categoryDtos, response.SpendingSummaryDTO{
			CategoryName: name,
			TotalSpent:   response.NewDecimal2(percentage),
		})
	}

	monthlySummaries, err := s.transactionRepo.GetMonthlySpendingSummary(ctx, accountID, 6)
	if err != nil {
		return nil, entities.NewAppError(entities.ErrTypeInternal, "failed to get monthly spending summary", err)
	}

	return &response.OverallSpendingSummaryDTO{
		ByCategory: categoryDtos,
		ByMonth:    s.fillMonthlySpendingGaps(now, monthlySummaries, 6),
	}, nil
}

func (s *transactionService) GetMonthlySpendingSummary(ctx context.Context, accountID uint) ([]response.MonthlySpendingDTO, error) {
	now := time.Now()
	summaries, err := s.transactionRepo.GetMonthlySpendingSummary(ctx, accountID, 6)
	if err != nil {
		return nil, entities.NewAppError(entities.ErrTypeInternal, "failed to get monthly spending summary", err)
	}

	return s.fillMonthlySpendingGaps(now, summaries, 6), nil
}

func (s *transactionService) fillMonthlySpendingGaps(now time.Time, summaries []entities.MonthlySpending, count int) []response.MonthlySpendingDTO {
	indexedMonthly := make(map[string]float64)
	for _, ms := range summaries {
		indexedMonthly[ms.Month] = ms.TotalSpent
	}

	dtos := make([]response.MonthlySpendingDTO, 0, count)
	for i := 0; i < count; i++ {
		// Use first day of month to avoid issues with varying month lengths
		date := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, -i, 0)
		monthLabel := date.Format("2006-01")

		totalSpent := indexedMonthly[monthLabel]

		dtos = append(dtos, response.MonthlySpendingDTO{
			Month:      monthLabel,
			TotalSpent: response.NewDecimal2(satangToTHB(totalSpent)),
		})
	}
	return dtos
}
