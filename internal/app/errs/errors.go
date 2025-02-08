package errs

import "errors"

// Юзер
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrEmailAlreadyUsed  = errors.New("email already used")
	ErrInvalidPwd        = errors.New("invalid password")
)

// Токен
var (
	ErrRefreshTokenRequired = errors.New("'refresh_token' is required")
	ErrTokenInvalid         = errors.New("token invalid")
	ErrInvalidTokenType     = errors.New("invalid token type (need REFRESH)")
	ErrMissingUsername      = errors.New("missing username in payload")
	ErrMissingUserID        = errors.New("missing user_id in payload")
	ErrTokenNotFound        = errors.New("token not found")
	ErrTokenExpired         = errors.New("token has expired")
	ErrValidateInRedis      = errors.New("failed to validate token in Redis")
	ErrFailedToRefresh      = errors.New("failed to refresh token")
	ErrFailedToSave         = errors.New("failed to save token")
)
