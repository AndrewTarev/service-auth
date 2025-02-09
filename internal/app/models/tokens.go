package models

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type InputRefresh struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
