package response

import "github.com/Narutchai01/solpay-core-service/internal/entities"

type TransactionDTO struct {
	ID              uint    `json:"id"`
	TransactionUUID string  `json:"transaction_uuid"`
	AccountID       uint    `json:"account_id" gorm:"not null;index"`
	CategoryID      string  `json:"category_id" gorm:"not null;"`
	TransactionType string  `json:"transaction_type" gorm:"not null;"`
	Status          string  `json:"status" gorm:"not null;default:'pending'"`
	THBAmount       float64 `json:"thb_amount" gorm:"not null;default:0"`
	USDTAmount      float64 `json:"usdt_amount" gorm:"not null;default:0"`
	Fee             float64 `json:"fee" gorm:"not null;default:0"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

func FormaterTransactionDTO(transaction *entities.TransactionEntity) *TransactionDTO {
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

func FormaterTransactionDTOS(transactions []entities.TransactionEntity) *[]TransactionDTO {
	var transactionDTOs []TransactionDTO
	for _, transaction := range transactions {
		transactionDTOs = append(transactionDTOs, *FormaterTransactionDTO(&transaction))
	}
	return &transactionDTOs
}
