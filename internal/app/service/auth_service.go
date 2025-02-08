package service

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	logger "github.com/sirupsen/logrus"

	"service-auth/internal/app/errs"
	"service-auth/internal/app/models"
	"service-auth/internal/infrastructure/auth_helper"
)

func (s *AuthService) CreateUser(ctx context.Context, user models.User) (int, error) {
	// Хэшируем пароль пользователя
	hashedPassword, err := auth_helper.HashPassword(user.Password)
	if err != nil {
		logger.Errorf(err.Error())
		return 0, err
	}

	user.Password = hashedPassword

	res, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}
	return res, nil
}

// GenerateTokens создает токены access, refresh, сохраняет refresh в redis
func (s *AuthService) GenerateTokens(ctx context.Context, username, password string) (models.Tokens, error) {
	user, err := s.repo.GetUser(ctx, username)
	if err != nil {
		return models.Tokens{}, err
	}

	checkPwd := auth_helper.CheckPasswordHash(password, user.Password)
	if checkPwd != nil {
		return models.Tokens{}, errs.ErrInvalidPwd
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(user.Username, int(user.ID), s.cfg.Auth.AccessTokenTTL)
	if err != nil {
		return models.Tokens{}, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.Username, int(user.ID), s.cfg.Auth.RefreshTokenTTL)
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
func (s *AuthService) RefreshTokens(ctx context.Context, oldRefreshToken string) (models.Tokens, error) {
	claims, err := s.ValidateRefreshToken(oldRefreshToken)
	if err != nil {
		return models.Tokens{}, err
	}

	username := claims.(jwt.MapClaims)["username"].(string)
	sub := claims.(jwt.MapClaims)["sub"].(float64)

	err = s.repo.FindTokenInRedis(ctx, oldRefreshToken)
	if err != nil {
		return models.Tokens{}, err
	}

	newAccessToken, err := s.jwtManager.GenerateAccessToken(username, int(sub), s.cfg.Auth.AccessTokenTTL)
	if err != nil {
		return models.Tokens{}, fmt.Errorf("error generating access token: %w", err)
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(username, int(sub), s.cfg.Auth.RefreshTokenTTL)
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

func (s *AuthService) ValidateRefreshToken(token string) (jwt.Claims, error) {
	claims, err := s.jwtManager.DecodeJWT(token)
	if err != nil {
		return nil, err
	}

	tokenType, ok := claims.(jwt.MapClaims)["token_type"].(string)
	if !ok || tokenType != auth_helper.RefreshToken {
		return nil, errs.ErrInvalidTokenType
	}

	_, ok = claims.(jwt.MapClaims)["username"].(string)
	if !ok {
		logger.Debugf("missing username in payload: %v", err)
		return nil, errs.ErrMissingUsername
	}

	_, ok = claims.(jwt.MapClaims)["sub"].(float64)
	if !ok {
		logger.Debugf("missing sub in payload: %v", err)
		return nil, errs.ErrMissingUserID
	}

	return claims, nil
}

// RevokeToken удаляет refresh токен из redis
func (s *AuthService) RevokeToken(ctx context.Context, token string) error {
	return s.repo.DeleteToken(ctx, token)
}
