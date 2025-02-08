package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"service-auth/internal/app/errs"
	"service-auth/internal/app/models"
	"service-auth/internal/app/service"
	"service-auth/internal/configs"
)

type Auth struct {
	services service.Service
	cfg      *configs.Config
}

func NewAuth(services service.Service, cfg *configs.Config) *Auth {
	return &Auth{
		services: services,
		cfg:      cfg,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Creates a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.User true "User details"
// @Success 201 {object} models.UserIdResponse "User ID"
// @Failure 400 {object} middleware.ValidationErrorResponse "Invalid input format or validation errors"
// @Failure 500 {object} middleware.ValidationErrorResponse "Internal server error"
// @Router /auth/register [post]
func (h *Auth) Register(ctx *gin.Context) {
	var input models.User

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(err)
		return
	}

	// Валидация входных данных
	if err := input.Validate(); err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			ctx.Error(validationErrs)
			return
		}
	}

	userID, err := h.services.CreateUser(ctx, input)
	if err != nil {
		ctx.Error(err) // Передаем бизнес-ошибку в middleware
		return
	}

	ctx.JSON(http.StatusCreated, models.UserIdResponse{
		UserId: userID})
}

// Login godoc
// @Summary User login
// @Description Logs in a user and returns access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.SignInInput true "Login credentials"
// @Success 200 {object} models.Tokens "Access and refresh tokens"
// @Failure 400 {object} middleware.ValidationErrorResponse "Invalid input format"
// @Failure 403 {object} middleware.ValidationErrorResponse "invalid username or password"
// @Router /auth/login [post]
func (h *Auth) Login(ctx *gin.Context) {
	var input models.SignInInput
	var tokens models.Tokens

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(err)
		return
	}

	tokens, err := h.services.GenerateTokens(ctx, input.Username, input.Password)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.SetCookie("access_token", tokens.AccessToken, int(h.cfg.Auth.AccessTokenTTL), "/", "", true, true)
	ctx.SetCookie("refresh_token", tokens.RefreshToken, int(h.cfg.Auth.RefreshTokenTTL), "/", "", true, true)

	ctx.JSON(http.StatusOK, tokens)
}

// Refresh godoc
// @Summary Refresh access token
// @Description Refreshes the access token using a valid refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.InputRefresh true "Refresh token"
// @Success 200 {object} models.Tokens "New access and refresh tokens"
// @Failure 400 {object} middleware.ValidationErrorResponse "Invalid input format"
// @Failure 401 {object} middleware.ValidationErrorResponse "Unauthorized or invalid refresh token"
// @Router /auth/refresh [post]
func (h *Auth) Refresh(ctx *gin.Context) {
	var tokens models.Tokens
	var input models.InputRefresh

	// Привязка и валидация входных данных
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(err)
		return
	}
	// Обновление токенов
	tokens, err := h.services.RefreshTokens(ctx, input.RefreshToken)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.SetCookie("access_token", tokens.AccessToken, int(h.cfg.Auth.AccessTokenTTL), "/", "", true, true)
	ctx.SetCookie("refresh_token", tokens.RefreshToken, int(h.cfg.Auth.RefreshTokenTTL), "/", "", true, true)

	ctx.JSON(http.StatusOK, tokens)
}

// Revoke godoc
// @Summary Revoke refresh token
// @Description Revokes the specified refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token query string true "Refresh token"
// @Success 200 {object} SuccessResponse "Token revoked successfully"
// @Failure 400 {object} middleware.ValidationErrorResponse "Invalid input format"
// @Failure 500 {object} middleware.ValidationErrorResponse "Internal server error"
// @Router /auth/revoke-token [delete]
func (h *Auth) Revoke(ctx *gin.Context) {
	token := ctx.Query("refresh_token")
	if token == "" {
		ctx.Error(errs.ErrRefreshTokenRequired)
		return
	}

	err := h.services.RevokeToken(ctx, token)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := SuccessResponse{Message: "Token revoked successfully"}

	ctx.JSON(http.StatusOK, response)
}
