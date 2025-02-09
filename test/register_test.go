package test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"service-auth/internal/app/errs"
	"service-auth/internal/app/models"
	"service-auth/internal/app/service"
	"service-auth/internal/app/service/mocks"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создание моков
	mockAuthService := mocks.NewMockAuthService(ctrl) // Используем сгенерированный мок

	services := service.Service{AuthService: mockAuthService}

	user := models.UserInput{
		Username: "testuser",
		Password: "password123",
		Email:    "testuser@gmail.com",
	}
	uuidObj, err := uuid.Parse("3410a4f3-a0b0-432a-89ba-9ea9ba48c6c1")
	// Настройка ожидания для метода CreateUser
	mockAuthService.EXPECT().CreateUser(gomock.Any(), user).Return(uuidObj, nil)

	// Test CreateUser
	id, err := services.CreateUser(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, uuidObj, id)
}

// Тест: Ошибка, если пользователь уже существует
func TestCreateUser_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthService(ctrl)

	svc := service.Service{AuthService: mockAuthService}

	user := models.UserInput{
		Username: "existing_user",
		Password: "password123",
		Email:    "test@example.com",
	}

	// Симулируем ошибку существующего пользователя
	mockAuthService.EXPECT().CreateUser(gomock.Any(), user).Return(uuid.UUID{}, errs.ErrUserAlreadyExists)

	_, err := svc.CreateUser(context.Background(), user)

	assert.Error(t, err)
	assert.Equal(t, errs.ErrUserAlreadyExists, err)
}

// Тест: Ошибка, если email уже используется
func TestCreateUser_EmailAlreadyUsed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthService(ctrl)

	svc := service.Service{AuthService: mockAuthService}

	user := models.UserInput{
		Username: "newuser",
		Password: "password123",
		Email:    "used_email@example.com",
	}

	// Симулируем ошибку email-а
	mockAuthService.EXPECT().CreateUser(gomock.Any(), user).Return(uuid.UUID{}, errs.ErrEmailAlreadyUsed)

	_, err := svc.CreateUser(context.Background(), user)

	assert.Error(t, err)
	assert.Equal(t, errs.ErrEmailAlreadyUsed, err)
}

func TestGenerateTokens(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthService(ctrl)

	svc := service.Service{AuthService: mockAuthService}

	username := "testuser"
	password := "password123"

	mockAuthService.EXPECT().GenerateTokens(gomock.Any(), username, password).
		Return(models.Tokens{
			AccessToken:  "accessToken",
			RefreshToken: "refreshToken",
		}, nil)

	tokens, err := svc.GenerateTokens(context.Background(), username, password)
	assert.NoError(t, err)
	assert.Equal(t, "accessToken", tokens.AccessToken)
	assert.Equal(t, "refreshToken", tokens.RefreshToken)
}

func TestUserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthService := mocks.NewMockAuthService(ctrl)
	svc := service.Service{AuthService: mockAuthService}
	username := "testuser"
	password := "password123"

	mockAuthService.EXPECT().GenerateTokens(gomock.Any(), username, password).Return(models.Tokens{}, errs.ErrUserNotFound)
	_, err := svc.GenerateTokens(context.Background(), username, password)
	assert.Error(t, err)
	assert.Equal(t, errs.ErrUserNotFound, err)
}

func TestInvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthService := mocks.NewMockAuthService(ctrl)
	svc := service.Service{AuthService: mockAuthService}

	username := "testuser"
	password := "password123"

	mockAuthService.EXPECT().GenerateTokens(gomock.Any(), username, password).Return(models.Tokens{}, errs.ErrInvalidPwd)
	_, err := svc.GenerateTokens(context.Background(), username, password)

	assert.Error(t, err)
	assert.Equal(t, errs.ErrInvalidPwd, err)
}

func TestRefreshTokens(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthService(ctrl)

	svc := service.Service{AuthService: mockAuthService}

	oldToken := "oldRefreshToken"

	mockAuthService.EXPECT().RefreshTokens(gomock.Any(), oldToken).
		Return(models.Tokens{
			AccessToken:  "newAccessToken",
			RefreshToken: "newRefreshToken",
		}, nil)

	tokens, err := svc.RefreshTokens(context.Background(), oldToken)
	assert.NoError(t, err)
	assert.Equal(t, "newAccessToken", tokens.AccessToken)
	assert.Equal(t, "newRefreshToken", tokens.RefreshToken)
}

func TestRefreshTokens_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthService := mocks.NewMockAuthService(ctrl)
	svc := service.Service{AuthService: mockAuthService}

	oldToken := "oldRefreshToken"
	mockAuthService.EXPECT().RefreshTokens(gomock.Any(), oldToken).Return(models.Tokens{}, errs.ErrTokenInvalid)
	_, err := svc.RefreshTokens(context.Background(), oldToken)
	assert.Error(t, err)
	assert.Equal(t, errs.ErrTokenInvalid, err)
}

func TestInvalidTokenType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthService := mocks.NewMockAuthService(ctrl)
	svc := service.Service{AuthService: mockAuthService}
	oldToken := "oldRefreshToken"
	mockAuthService.EXPECT().RefreshTokens(gomock.Any(), oldToken).Return(models.Tokens{}, errs.ErrInvalidTokenType)
	_, err := svc.RefreshTokens(context.Background(), oldToken)
	assert.Error(t, err)
	assert.Equal(t, errs.ErrInvalidTokenType, err)
}

func TestTokenExpired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthService := mocks.NewMockAuthService(ctrl)
	svc := service.Service{AuthService: mockAuthService}
	oldToken := "oldRefreshToken"
	mockAuthService.EXPECT().RefreshTokens(gomock.Any(), oldToken).Return(models.Tokens{}, errs.ErrTokenExpired)
	_, err := svc.RefreshTokens(context.Background(), oldToken)
	assert.Error(t, err)
	assert.Equal(t, errs.ErrTokenExpired, err)
}

func TestRevokeToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthService(ctrl)

	svc := service.Service{AuthService: mockAuthService}

	token := "refreshToken"

	mockAuthService.EXPECT().RevokeToken(gomock.Any(), token).Return(nil)

	err := svc.RevokeToken(context.Background(), token)
	assert.NoError(t, err)
}

func TestRevokeToken_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthService := mocks.NewMockAuthService(ctrl)
	svc := service.Service{AuthService: mockAuthService}
	token := "refreshToken"
	mockAuthService.EXPECT().RevokeToken(gomock.Any(), token).Return(errs.ErrTokenInvalid)
	err := svc.RevokeToken(context.Background(), token)
	assert.Error(t, err)
	assert.Equal(t, errs.ErrTokenInvalid, err)
}

func TestRevokeToken_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthService := mocks.NewMockAuthService(ctrl)
	svc := service.Service{AuthService: mockAuthService}
	token := "refreshToken"
	mockAuthService.EXPECT().RevokeToken(gomock.Any(), token).Return(errs.ErrTokenNotFound)
	err := svc.RevokeToken(context.Background(), token)
	assert.Error(t, err)
	assert.Equal(t, errs.ErrTokenNotFound, err)
}
