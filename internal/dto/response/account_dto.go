package response

import "github.com/Narutchai01/solpay-core-service/internal/entities"

// AccountDTO is the API representation of an account.
type AccountDTO struct {
	ID            uint          `json:"id"`
	PublicAddress string        `json:"public_address"`
	IsKYCVerified bool          `json:"is_kyc_verified"`
	CreatedAt     string        `json:"created_at"`
	UpdatedAt     string        `json:"updated_at"`
	User          *UserResponse `json:"user"`
}

// FormatAccountDTO converts an AccountEntity to its API representation.
func FormatAccountDTO(account *entities.AccountEntity) *AccountDTO {
	return &AccountDTO{
		ID:            account.ID,
		PublicAddress: account.PublicAddress,
		IsKYCVerified: isKYCVerified(account),
		CreatedAt:     account.CreatedAt.String(),
		UpdatedAt:     account.UpdatedAt.String(),
		User:          FormatUserResponse(account.User),
	}
}

func isKYCVerified(account *entities.AccountEntity) bool {
	if account.KycToken == nil {
		return false
	}

	return *account.KycToken != ""
}

// FormatAccountDTOs converts a slice of AccountEntity to a slice of AccountDTO.
func FormatAccountDTOs(accounts []entities.AccountEntity) []AccountDTO {
	dtos := make([]AccountDTO, 0, len(accounts))
	for _, account := range accounts {
		dtos = append(dtos, *FormatAccountDTO(&account))
	}
	return dtos
}
