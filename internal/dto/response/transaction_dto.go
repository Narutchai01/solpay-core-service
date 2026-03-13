package response

import "github.com/Narutchai01/solpay-core-service/internal/entities"

// TransactionDTO is the API representation of a transaction.
type TransactionDTO struct {
	ID              uint    `json:"id"`
	TransactionUUID string  `json:"transaction_uuid"`
	AccountID       uint    `json:"account_id"`
	CategoryID      string  `json:"category_id"`
	TransactionType string  `json:"transaction_type"`
	Status          string  `json:"status"`
	THBAmount       float64 `json:"thb_amount"`
	USDTAmount      float64 `json:"usdt_amount"`
	Fee             float64 `json:"fee"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// FormatTransactionDTO converts a TransactionEntity to its API representation.
func FormatTransactionDTO(transaction *entities.TransactionEntity) *TransactionDTO {
	return &TransactionDTO{
		ID:              transaction.ID,
		TransactionUUID: transaction.TransactionUUID.String(),
		AccountID:       transaction.AccountID,
		CategoryID:      transaction.CategoryID,
		TransactionType: transaction.TransactionType,
		Status:          transaction.Status,
		THBAmount:       transaction.THBAmount,
		USDTAmount:      transaction.USDTAmount,
		Fee:             transaction.Fee,
		CreatedAt:       transaction.CreatedAt.String(),
		UpdatedAt:       transaction.UpdatedAt.String(),
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
