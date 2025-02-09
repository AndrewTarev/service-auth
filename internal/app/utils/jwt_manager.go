package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
func (j *JWTManager) GenerateAccessToken(username, role string, id uuid.UUID, accessTTL time.Duration) (string, error) {
	logger.Debug("Generating access token for user: ", username)

	claims := jwt.MapClaims{
		"token_type": AccessToken,
		"sub":        id,
		"username":   username,
		"role":       role,
		"exp":        time.Now().Add(accessTTL).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.jwtSecret)
}

// GenerateRefreshToken создает токен обновления
func (j *JWTManager) GenerateRefreshToken(username, role string, id uuid.UUID, refreshTTL time.Duration) (string, error) {
	logger.Debug("Generating refresh token")

	claims := jwt.MapClaims{
		"token_type": RefreshToken,
		"sub":        id,
		"username":   username,
		"role":       role,
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
		logger.Debugf("Error parsing token: %v", err)
		return nil, errs.ErrTokenInvalid
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errs.ErrTokenInvalid
}
