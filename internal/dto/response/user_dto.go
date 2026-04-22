package response

type UserResponse struct {
	ID           uint   `json:"id"`
	IDCard       string `json:"id_card"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	BirthDate    string `json:"birth_date"`
	Status       string `json:"status"`
	ExpireDate   string `json:"expire_date"`
	FrontCardURL string `json:"front_card_url"`
	BackCardURL  string `json:"back_card_url"`
}
