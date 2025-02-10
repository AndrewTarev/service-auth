package utils

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"service-auth/internal/app/errs"
)

const (
	AccessToken  = "access"
	RefreshToken = "refresh"
)

// JWTManager управляет генерацией токенов
type JWTManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewJWTManager загружает RSA-ключи и создает JWT-менеджер
func NewJWTManager(privateKeyPath, publicKeyPath string) (*JWTManager, error) {
	// Загружаем приватный ключ
	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		logger.WithError(err).Error("failed to read private key")
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		logger.WithError(err).Error("failed to parse private key")
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Загружаем публичный ключ
	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		logger.WithError(err).Error("failed to read public key")
		return nil, fmt.Errorf("failed to read public key: %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		logger.WithError(err).Error("failed to parse public key")
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return &JWTManager{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
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

// GenerateAccessToken создает токен доступа (подписан приватным ключом)
func (j *JWTManager) GenerateAccessToken(username, role string, id uuid.UUID, accessTTL time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"token_type": "access",
		"sub":        id.String(),
		"username":   username,
		"role":       role,
		"exp":        time.Now().Add(accessTTL).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(j.privateKey)
}

// GenerateRefreshToken создает refresh токен (подписан приватным ключом)
func (j *JWTManager) GenerateRefreshToken(username, role string, id uuid.UUID, refreshTTL time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"token_type": "refresh",
		"sub":        id.String(),
		"username":   username,
		"role":       role,
		"exp":        time.Now().Add(refreshTTL).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(j.privateKey)
}

// DecodeJWT парсит токен и проверяет его подпись публичным ключом
func (j *JWTManager) DecodeJWT(tokenString string) (jwt.Claims, error) {
	logger.Debug("Parsing token")

	// Разбираем и проверяем подпись токена
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что используется алгоритм RS256
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.publicKey, nil // Возвращаем публичный ключ для проверки подписи
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			logger.Debugf("Error: token is expired: %v", err)
			return nil, errs.ErrTokenExpired
		}
		logger.Debugf("Error parsing token: %v", err)
		return nil, errs.ErrTokenInvalid
	}

	// Проверяем валидность токена
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errs.ErrTokenInvalid
}
