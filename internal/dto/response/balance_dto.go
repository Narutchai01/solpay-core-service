package response

import (
	"github.com/Narutchai01/solpay-core-service/internal/entities"
)

type BalanceDTO struct {
	ID         uint    `json:"id"`
	AccountID  uint    `json:"account_id"`
	THBAmount  float64 `json:"thb_amount"`
	USDTAmount float64 `json:"usdt_amount"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

func FormaterBalanceDTO(balance *entities.BalanceEntity) *BalanceDTO {
	thbFloat := float64(balance.THBAmount) / 10000.0
	usdtFloat := float64(balance.USDTAmount) / 10000.0

	return &BalanceDTO{
		ID:         balance.ID,
		AccountID:  balance.AccountID,
		THBAmount:  thbFloat,
		USDTAmount: usdtFloat,
		CreatedAt:  balance.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  balance.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func FormaterBalanceDTOS(balances []entities.BalanceEntity) []BalanceDTO {
	balanceDTOs := make([]BalanceDTO, 0, len(balances))
	for _, balance := range balances {
		balanceDTOs = append(balanceDTOs, *FormaterBalanceDTO(&balance))
	}
	return balanceDTOs
}
