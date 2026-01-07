package response

import "github.com/Narutchai01/solpay-core-service/internal/entities"

type AccountDTO struct {
	ID            uint   `json:"id"`
	PublicAddress string `json:"public_address"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

func FormaterAccountDTO(account entities.AccountEntity) AccountDTO {
	return AccountDTO{
		ID:            account.ID,
		PublicAddress: account.PublicAddress,
		CreatedAt:     account.CreatedAt.String(),
		UpdatedAt:     account.UpdatedAt.String(),
	}
}

func FormaterAccountDTOS(accounts []entities.AccountEntity) []AccountDTO {
	var accountModels []AccountDTO
	for _, account := range accounts {
		accountModels = append(accountModels, FormaterAccountDTO(account))
	}
	return accountModels
}
