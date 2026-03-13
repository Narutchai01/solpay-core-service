package response

// TokenDTO represents the authentication token pair returned to clients.
type TokenDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
