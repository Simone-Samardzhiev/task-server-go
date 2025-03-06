package models

// TokenGroup struct holds information both about access and refresh token.
type TokenGroup struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewTokenGroup(accessToken string, refreshToken string) *TokenGroup {
	return &TokenGroup{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
