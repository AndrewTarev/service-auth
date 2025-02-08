package service

import (
	"context"

	"service-auth/internal/app/models"
	"service-auth/internal/app/repository"
	"service-auth/internal/app/utils"
	"service-auth/internal/configs"
)

type AuthService interface {
	CreateUser(ctx context.Context, user models.User) (int, error)
	GenerateTokens(ctx context.Context, username, password string) (models.Tokens, error)
	RefreshTokens(ctx context.Context, oldRefreshToken string) (models.Tokens, error)
	RevokeToken(ctx context.Context, token string) error
}

type Service struct {
	AuthService
}

func NewService(repo *repository.Repository, jwtManager *utils.JWTManager, cfg *configs.Config) Service {
	return Service{
		AuthService: NewAuth(repo, jwtManager, cfg),
	}
}
