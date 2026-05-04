package response

import (
	"math"
	"strconv"

	"github.com/Narutchai01/solpay-core-service/internal/entities"
)

// Decimal2 is encoded as a JSON number with exactly 2 decimal places.
type Decimal2 float64

func NewDecimal2(v float64) Decimal2 {
	return Decimal2(math.Round(v*100) / 100)
}

func (d Decimal2) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatFloat(float64(d), 'f', 2, 64)), nil
}

// Decimal6 is encoded as a JSON number with exactly 6 decimal places.
type Decimal6 float64

func NewDecimal6(v float64) Decimal6 {
	return Decimal6(math.Round(v*1000000) / 1000000)
}

func (d Decimal6) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatFloat(float64(d), 'f', 6, 64)), nil
}

// Decimal9 is encoded as a JSON number with exactly 9 decimal places.
type Decimal9 float64

func NewDecimal9(v float64) Decimal9 {
	return Decimal9(math.Round(v*1000000000) / 1000000000)
}

func (d Decimal9) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatFloat(float64(d), 'f', 9, 64)), nil
}

type TransactionOnChainDTO struct {
	TxHash    string `json:"tx_hash"`
	Signature string `json:"signature"`
}

type TransactionOffChainDTO struct {
	PromptPayID string  `json:"prompt_pay_id"`
	SlipURL     *string `json:"slip_url,omitempty"`
}

// TransactionDTO is the API representation of a transaction.
type TransactionDTO struct {
	ID                  uint                    `json:"id"`
	TransactionUUID     string                  `json:"transaction_uuid"`
	AccountID           uint                    `json:"account_id"`
	TransactionType     string                  `json:"transaction_type"`
	Status              string                  `json:"status"`
	Category            *entities.Category      `json:"category"`
	THBAmount           Decimal2                `json:"thb_amount"`
	USDTAmount          Decimal6                `json:"usdt_amount"`
	SOLAmount           Decimal9                `json:"sol_amount"`
	Fee                 float64                 `json:"fee"`
	TransactionOnChain  *TransactionOnChainDTO  `json:"transaction_on_chain,omitempty"`
	TransactionOffChain *TransactionOffChainDTO `json:"transaction_off_chain,omitempty"`
	CreatedAt           string                  `json:"created_at"`
	AccountDTO          *AccountDTO             `json:"account,omitempty"`
	UpdatedAt           string                  `json:"updated_at"`
}

// FormatTransactionDTO converts a TransactionEntity to its API representation.
func FormatTransactionDTO(transaction *entities.TransactionEntity) *TransactionDTO {
	var onChain *TransactionOnChainDTO
	if transaction.TransactionOnChain != nil {
		onChain = &TransactionOnChainDTO{
			TxHash:    transaction.TransactionOnChain.TxHash,
			Signature: transaction.TransactionOnChain.Signature,
		}
	}

	var offChain *TransactionOffChainDTO
	if transaction.TransactionOffChain != nil {
		offChain = &TransactionOffChainDTO{
			PromptPayID: transaction.TransactionOffChain.PromptPayID,
			SlipURL:     transaction.TransactionOffChain.SlipURL,
		}
	}

	var accountDTO *AccountDTO
	if transaction.Account != nil {
		accountDTO = FormatAccountDTO(transaction.Account)
	}

	return &TransactionDTO{
		ID:                  transaction.ID,
		TransactionUUID:     transaction.TransactionUUID.String(),
		AccountID:           transaction.AccountID,
		TransactionType:     transaction.TransactionType,
		Status:              transaction.Status,
		Category:            transaction.Category,
		THBAmount:           NewDecimal2(transaction.THBAmount / 100),
		USDTAmount:          NewDecimal6(transaction.USDTAmount),
		SOLAmount:           NewDecimal9(transaction.SOLAmount),
		Fee:                 transaction.Fee,
		TransactionOnChain:  onChain,
		TransactionOffChain: offChain,
		CreatedAt:           transaction.CreatedAt.String(),
		AccountDTO:          accountDTO,
		UpdatedAt:           transaction.UpdatedAt.String(),
	}
}

// FormatTransactionDTOs converts a slice of TransactionEntity to a slice of TransactionDTO.
func FormatTransactionDTOs(transactions []entities.TransactionEntity) []TransactionDTO {
	dtos := make([]TransactionDTO, 0, len(transactions))
	for _, tx := range transactions {
		dtos = append(dtos, *FormatTransactionDTO(&tx))
	}
	return dtos
}

type TransactionChartData struct {
	Date     string   `json:"date"`
	Label    string   `json:"label"`
	Deposit  float64  `json:"deposit"`
	Withdraw Decimal2 `json:"withdraw"`
	Fee      float64  `json:"fee"`
}

type TransactionChartSummary struct {
	TotalDeposit        float64                `json:"totalDeposit"`
	TotalWithdraw       Decimal2               `json:"totalWithdraw"`
	TotalFee            float64                `json:"totalFee"`
	TotalCompletedCount int                    `json:"totalCompletedCount"`
	ChartData           []TransactionChartData `json:"chartData"`
}

type SpendingSummaryDTO struct {
	CategoryName string   `json:"category_name"`
	TotalSpent   Decimal2 `json:"total_spent"`
}

type MonthlySpendingDTO struct {
	Month      string   `json:"month"`
	TotalSpent Decimal2 `json:"total_spent"`
}

type OverallSpendingSummaryDTO struct {
	ByCategory []SpendingSummaryDTO `json:"by_category"`
	ByMonth    []MonthlySpendingDTO `json:"by_month"`
}
