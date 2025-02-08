package models

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type InputRefresh struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AccessToken struct {
	TokenType string `json:"token_type"`
	Subject   int    `json:"sub"`
	Username  string `json:"username"`
	Expires   int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
}
