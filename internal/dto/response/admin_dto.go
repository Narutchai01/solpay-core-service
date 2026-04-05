package response

import "github.com/Narutchai01/solpay-core-service/internal/entities"

type AdminDTO struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func FormatAdminDTO(admin *entities.AdminEntity) *AdminDTO {
	return &AdminDTO{
		ID:        admin.ID.String(),
		Username:  admin.Username,
		CreatedAt: admin.CreatedAt.String(),
		UpdatedAt: admin.UpdatedAt.String(),
	}
}
