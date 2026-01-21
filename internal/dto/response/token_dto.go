package response

type TokenDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func FormaterTokenResponse(token *TokenDTO) *TokenDTO {
	return &TokenDTO{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
}
