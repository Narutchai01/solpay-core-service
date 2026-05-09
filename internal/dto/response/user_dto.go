package response

import "github.com/Narutchai01/solpay-core-service/internal/entities"

type UserResponse struct {
	ID           uint   `json:"id"`
	IDCard       string `json:"id_card"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	BirthDate    string `json:"birth_date"`
	Status       string `json:"status"`
	ExpireDate   string `json:"expire_date"`
	FaceURL      string `json:"face_url"`
	FrontCardURL string `json:"front_card_url"`
	BackCardURL  string `json:"back_card_url"`
}

func FormatUserResponse(user *entities.User) *UserResponse {
	if user == nil {
		return nil
	}
	return &UserResponse{
		ID:           user.ID,
		IDCard:       user.IDCard,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		BirthDate:    user.BirthDate,
		Status:       user.Status,
		ExpireDate:   user.ExpireDate,
		FaceURL:      user.FaceURL,
		FrontCardURL: user.FrontCardURL,
		BackCardURL:  user.BackCardURL,
	}
}

type UserListResponse struct {
	Rows       []*UserResponse `json:"rows"`
	TotalCount int             `json:"total_count"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
}
