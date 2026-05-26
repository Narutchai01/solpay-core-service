package services_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/core/services"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockTransactionRepository
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) CreateTransaction(ctx context.Context, tx *entities.TransactionEntity) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *MockTransactionRepository) CreateTransactionOnChain(ctx context.Context, tx *entities.TransactionOnChain) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *MockTransactionRepository) UpdateTransactionOnChain(ctx context.Context, tx *entities.TransactionOnChain) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *MockTransactionRepository) CreateTransactionOffChain(ctx context.Context, tx *entities.TransactionOffChain) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *MockTransactionRepository) UpdateTransactionOffChain(ctx context.Context, tx *entities.TransactionOffChain) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *MockTransactionRepository) UpdateTransactionStatus(ctx context.Context, txID string, status string) error {
	args := m.Called(ctx, txID, status)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetTransactionByID(id int) (*entities.TransactionEntity, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.TransactionEntity), args.Error(1)
}

func (m *MockTransactionRepository) GetTransactionByAccountID(accountID int) ([]entities.TransactionEntity, error) {
	args := m.Called(accountID)
	return args.Get(0).([]entities.TransactionEntity), args.Error(1)
}

func (m *MockTransactionRepository) GetTransactionByUUID(txUUID uuid.UUID) (*entities.TransactionEntity, error) {
	args := m.Called(txUUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.TransactionEntity), args.Error(1)
}

func (m *MockTransactionRepository) CountTransactions(query request.TransactionQuery, accountID *uint) (int64, error) {
	args := m.Called(query, accountID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTransactionRepository) GetTransactions(query request.TransactionQuery, accountID *uint) ([]entities.TransactionEntity, error) {
	args := m.Called(query, accountID)
	return args.Get(0).([]entities.TransactionEntity), args.Error(1)
}

func (m *MockTransactionRepository) QueryTransactionSummary(ctx context.Context, month, year int) ([]entities.TransactionSummary, error) {
	args := m.Called(ctx, month, year)
	return args.Get(0).([]entities.TransactionSummary), args.Error(1)
}

func (m *MockTransactionRepository) GetSpendingSummary(ctx context.Context, accountID uint, month, year int) ([]entities.SpendingSummary, error) {
	args := m.Called(ctx, accountID, month, year)
	return args.Get(0).([]entities.SpendingSummary), args.Error(1)
}

func (m *MockTransactionRepository) GetMonthlySpendingSummary(ctx context.Context, accountID uint, limit int) ([]entities.MonthlySpending, error) {
	args := m.Called(ctx, accountID, limit)
	return args.Get(0).([]entities.MonthlySpending), args.Error(1)
}

// MockUnitOfWork
type MockUnitOfWork struct {
	mock.Mock
}

func (m *MockUnitOfWork) Execute(ctx context.Context, fn func(ctx context.Context) (any, error)) (any, error) {
	return fn(ctx)
}

// MockPublisher
type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) Publish(queueName string, body []byte) error {
	args := m.Called(queueName, body)
	return args.Error(0)
}

func (m *MockPublisher) PublishTransactionMessage(queueName string, tx *entities.TransactionEntity, status string) {
	m.Called(queueName, tx, status)
}

// MockWebSocketPort
type MockWebSocketPort struct {
	mock.Mock
}

func (m *MockWebSocketPort) NotifyTransactionStatus(txID string, payload []byte) {
	m.Called(txID, payload)
}

type TransactionServiceTestSuite struct {
	suite.Suite
	repo      *MockTransactionRepository
	uow       *MockUnitOfWork
	publisher *MockPublisher
	wsHub     *MockWebSocketPort
	service   services.TransactionService
}

func (suite *TransactionServiceTestSuite) SetupTest() {
	suite.repo = new(MockTransactionRepository)
	suite.uow = new(MockUnitOfWork)
	suite.publisher = new(MockPublisher)
	suite.wsHub = new(MockWebSocketPort)
	cfg := &config.Config{
		SOLANA_WORK_QUEUE:              "test-solana-queue",
		BALANCE_QUEUE:                  "test-balance-queue",
		PAYMENT_QUEUE:                  "test-payment-queue",
		TRANSACTION_ORCHESTRATOR_QUEUE: "test-orch-queue",
	}
	suite.service = services.NewTransactionService(suite.repo, suite.uow, suite.publisher, suite.wsHub, cfg)
}

func TestTransactionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionServiceTestSuite))
}

func (suite *TransactionServiceTestSuite) TestHandleTransactionUpdate_FinalizeFailedOnWorkflowError() {
	txUUID := uuid.New()
	tx := &entities.TransactionEntity{
		TransactionUUID: txUUID,
		TransactionType: string(entities.TOPUP),
		Status:          string(entities.StatusPending),
		TransactionOnChain: &entities.TransactionOnChain{
			TxHash: "test-hash",
		},
	}

	event := request.TransactionMessage{
		TxID:   txUUID.String(),
		Status: string(entities.StatusSolanaSubmitted),
	}
	msg, _ := json.Marshal(event)

	// Mock repository behavior
	suite.repo.On("GetTransactionByUUID", txUUID).Return(tx, nil)
	suite.repo.On("UpdateTransactionStatus", mock.Anything, txUUID.String(), string(entities.StatusSolanaSubmitted)).Return(nil)

	// Force a failure in the workflow (e.g. publisher failure)
	suite.publisher.On("Publish", mock.Anything, mock.Anything).Return(errors.New("publisher error"))

	// Expected finalized update
	suite.repo.On("UpdateTransactionStatus", mock.Anything, txUUID.String(), string(entities.StatusFailed)).Return(nil)
	suite.wsHub.On("NotifyTransactionStatus", txUUID.String(), mock.Anything).Return()

	err := suite.service.HandleTransactionUpdate(context.Background(), msg)

	suite.Error(err)
	suite.Contains(err.Error(), "publisher error")
	suite.repo.AssertCalled(suite.T(), "UpdateTransactionStatus", mock.Anything, txUUID.String(), string(entities.StatusFailed))
}

func (suite *TransactionServiceTestSuite) TestHandleTransactionUpdate_IncomingFailedStatus() {
	txUUID := uuid.New()
	tx := &entities.TransactionEntity{
		TransactionUUID: txUUID,
		TransactionType: string(entities.TOPUP),
		Status:          string(entities.StatusSolanaSubmitted),
	}

	event := request.TransactionMessage{
		TxID:   txUUID.String(),
		Status: string(entities.StatusSolanaFailed),
	}
	msg, _ := json.Marshal(event)

	suite.repo.On("GetTransactionByUUID", txUUID).Return(tx, nil)
	// Initial update to SOLANA_FAILED
	suite.repo.On("UpdateTransactionStatus", mock.Anything, txUUID.String(), string(entities.StatusSolanaFailed)).Return(nil)
	// Finalize update to FAILED (mapped from SOLANA_FAILED in topupWorkflow)
	suite.repo.On("UpdateTransactionStatus", mock.Anything, txUUID.String(), string(entities.StatusFailed)).Return(nil)
	suite.wsHub.On("NotifyTransactionStatus", txUUID.String(), mock.Anything).Return()

	err := suite.service.HandleTransactionUpdate(context.Background(), msg)

	suite.NoError(err)
	suite.repo.AssertCalled(suite.T(), "UpdateTransactionStatus", mock.Anything, txUUID.String(), string(entities.StatusFailed))
}
