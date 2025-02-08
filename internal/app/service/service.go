package service

import (
	"context"

	"service-auth/internal/app/models"
	"service-auth/internal/app/repository"
	"service-auth/internal/configs"
	"service-auth/internal/infrastructure/auth_helper"
)

type Service interface {
	CreateUser(ctx context.Context, user models.User) (int, error)
	GenerateTokens(ctx context.Context, username, password string) (models.Tokens, error)
	RefreshTokens(ctx context.Context, oldRefreshToken string) (models.Tokens, error)
	RevokeToken(ctx context.Context, token string) error
}

type AuthService struct {
	repo       *repository.Repository
	jwtManager *auth_helper.JWTManager
	cfg        *configs.Config
}

func NewService(repo *repository.Repository, jwtManager *auth_helper.JWTManager, cfg *configs.Config) Service {
	return &AuthService{repo: repo, jwtManager: jwtManager, cfg: cfg}
}
