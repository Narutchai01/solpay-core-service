package response

import (
	"github.com/Narutchai01/solpay-core-service/internal/entities"
)

// BalanceDTO is the API representation of a balance.
type BalanceDTO struct {
	ID         uint    `json:"id"`
	AccountID  uint    `json:"account_id"`
	THBAmount  float64 `json:"thb_amount"`
	USDTAmount float64 `json:"usdt_amount"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

// FormatBalanceDTO converts a BalanceEntity to its API representation.
func FormatBalanceDTO(balance *entities.BalanceEntity) *BalanceDTO {
	const divisor = 10000.0

	return &BalanceDTO{
		ID:        balance.ID,
		AccountID: balance.AccountID,
		THBAmount: float64(balance.THBAmount) / divisor,
		CreatedAt: balance.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: balance.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// FormatBalanceDTOs converts a slice of BalanceEntity to a slice of BalanceDTO.
func FormatBalanceDTOs(balances []entities.BalanceEntity) []BalanceDTO {
	dtos := make([]BalanceDTO, 0, len(balances))
	for _, balance := range balances {
		dtos = append(dtos, *FormatBalanceDTO(&balance))
	}
	return dtos
}
