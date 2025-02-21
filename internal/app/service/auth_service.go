package service

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	logger "github.com/sirupsen/logrus"

	"service-auth/internal/app/errs"
	"service-auth/internal/app/models"
	"service-auth/internal/app/repository"
	"service-auth/internal/app/utils"
	"service-auth/internal/configs"
)

type Auth struct {
	repo       *repository.Repository
	jwtManager *utils.JWTManager
	cfg        *configs.Config
}

func NewAuth(repo *repository.Repository, jwtManager *utils.JWTManager, cfg *configs.Config) *Auth {
	return &Auth{repo: repo, jwtManager: jwtManager, cfg: cfg}
}

func (s *Auth) CreateUser(ctx context.Context, user models.UserInput) (uuid.UUID, error) {
	// Хэшируем пароль пользователя
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		logger.Errorf(err.Error())
		return uuid.UUID{}, err
	}

	user.Password = hashedPassword

	res, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return uuid.UUID{}, err
	}
	return res, nil
}

// GenerateTokens создает токены access, refresh, сохраняет refresh в redis
func (s *Auth) GenerateTokens(ctx context.Context, username, password string) (models.Tokens, error) {
	user, err := s.repo.GetUser(ctx, username)
	if err != nil {
		return models.Tokens{}, err
	}

	checkPwd := utils.CheckPasswordHash(password, user.Password)
	if checkPwd != nil {
		return models.Tokens{}, errs.ErrInvalidPwd
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(user.Username, user.Role, user.ID, s.cfg.Auth.AccessTokenTTL)
	if err != nil {
		return models.Tokens{}, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.Username, user.Role, user.ID, s.cfg.Auth.RefreshTokenTTL)
	if err != nil {
		return models.Tokens{}, err
	}

	// сохраняем refresh в redis
	err = s.repo.SaveToken(ctx, refreshToken, s.cfg.Auth.RefreshTokenTTL)
	if err != nil {
		return models.Tokens{}, err
	}

	return models.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshTokens обновляет access и refresh токен доступа
func (s *Auth) RefreshTokens(ctx context.Context, oldRefreshToken string) (models.Tokens, error) {
	claims, err := s.ValidateRefreshToken(oldRefreshToken)
	if err != nil {
		return models.Tokens{}, err
	}

	username := claims.(jwt.MapClaims)["username"].(string)
	role := claims.(jwt.MapClaims)["role"].(string)
	sub := claims.(jwt.MapClaims)["sub"].(string)

	uuidObj, err := uuid.Parse(sub)
	if err != nil {
		logger.Errorf("error parse UUID: %s", err.Error())
		return models.Tokens{}, errs.ErrParseUUID
	}

	err = s.repo.FindTokenInRedis(ctx, oldRefreshToken)
	if err != nil {
		return models.Tokens{}, err
	}

	newAccessToken, err := s.jwtManager.GenerateAccessToken(username, role, uuidObj, s.cfg.Auth.AccessTokenTTL)
	if err != nil {
		return models.Tokens{}, fmt.Errorf("error generating access token: %w", err)
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(username, role, uuidObj, s.cfg.Auth.RefreshTokenTTL)
	if err != nil {
		return models.Tokens{}, fmt.Errorf("error generating refresh token: %w", err)
	}

	logger.Debug("Access tokens refreshed successfully for user: ", username)

	err = s.repo.UpdateRefreshTokenInRedis(ctx, oldRefreshToken, newRefreshToken, s.cfg.Auth.RefreshTokenTTL)
	if err != nil {
		return models.Tokens{}, fmt.Errorf("failed to update refresh token: %w", err)
	}
	return models.Tokens{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *Auth) ValidateRefreshToken(token string) (jwt.Claims, error) {
	claims, err := s.jwtManager.DecodeJWT(token)
	if err != nil {
		return nil, errs.ErrTokenInvalid
	}

	tokenType, ok := claims.(jwt.MapClaims)["token_type"].(string)
	if !ok || tokenType != utils.RefreshToken {
		return nil, errs.ErrInvalidTokenType
	}

	_, ok = claims.(jwt.MapClaims)["username"].(string)
	if !ok {
		logger.Debugf("missing username in payload: %v", err)
		return nil, errs.ErrMissingUsername
	}

	_, ok = claims.(jwt.MapClaims)["sub"].(string)
	if !ok {
		logger.Debugf("missing sub in payload: %v", err)
		return nil, errs.ErrMissingUserID
	}

	return claims, nil
}

// RevokeToken удаляет refresh токен из redis
func (s *Auth) RevokeToken(ctx context.Context, token string) error {
	return s.repo.DeleteToken(ctx, token)
}
