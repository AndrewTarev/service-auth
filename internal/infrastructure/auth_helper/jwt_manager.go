package auth_helper

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"service-auth/internal/app/errs"
	"service-auth/internal/configs"
)

const (
	AccessToken  = "access"
	RefreshToken = "refresh"
)

type JWTManager struct {
	jwtSecret []byte
}

func NewJWTManager(cfg *configs.Config) *JWTManager {
	return &JWTManager{
		jwtSecret: []byte(cfg.Auth.SigningKey),
	}
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}

	return string(hashed), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// GenerateAccessToken создает токен доступа с информацией о пользователе
func (j *JWTManager) GenerateAccessToken(username string, id int, accessTTL time.Duration) (string, error) {
	logger.Debug("Generating access token for user: ", username)

	claims := jwt.MapClaims{
		"token_type": AccessToken,
		"sub":        id,
		"username":   username,
		"exp":        time.Now().Add(accessTTL).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.jwtSecret)
}

// GenerateRefreshToken создает токен обновления
func (j *JWTManager) GenerateRefreshToken(username string, id int, refreshTTL time.Duration) (string, error) {
	logger.Debug("Generating refresh token")

	claims := jwt.MapClaims{
		"token_type": RefreshToken,
		"sub":        id,
		"username":   username,
		"exp":        time.Now().Add(refreshTTL).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.jwtSecret)
}

// DecodeJWT парсит токен и возвращает его клеймы
func (j *JWTManager) DecodeJWT(tokenString string) (jwt.Claims, error) {
	logger.Debug("Parse token")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			logger.Debugf("Error: token is expired: %v", err)
			return nil, errs.ErrTokenExpired
		}
		log.Printf("Error parsing token: %v", err)
		return nil, errs.ErrTokenInvalid
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errs.ErrTokenInvalid
}

// RefreshTokens обновляет токены доступа и обновления
func (j *JWTManager) RefreshTokens(refreshToken string, accessTTL, refreshTTL time.Duration) (string, string, error) {
	logger.Debug("Refreshing tokens")

	claims, err := j.DecodeJWT(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %w", err)
	}

	tokenType, ok := claims.(jwt.MapClaims)["token_type"].(string)
	if !ok || tokenType != RefreshToken {
		return "", "", errs.ErrInvalidTokenType
	}

	username, ok := claims.(jwt.MapClaims)["username"].(string)
	if !ok {
		return "", "", errs.ErrMissingUsername
	}

	id, ok := claims.(jwt.MapClaims)["sub"].(float64)
	if !ok {
		return "", "", errs.ErrMissingUserID
	}

	newAccessToken, err := j.GenerateAccessToken(username, int(id), accessTTL)
	if err != nil {
		return "", "", fmt.Errorf("error generating access token: %w", err)
	}

	newRefreshToken, err := j.GenerateRefreshToken(username, int(id), refreshTTL)
	if err != nil {
		return "", "", fmt.Errorf("error generating refresh token: %w", err)
	}

	logger.Debug("Access tokens refreshed successfully for user: ", username)
	return newAccessToken, newRefreshToken, nil
}
