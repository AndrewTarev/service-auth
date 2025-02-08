package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	logger "github.com/sirupsen/logrus"

	"service-auth/internal/app/errs"
)

// ValidationErrorResponse структура для JSON-ответа
type ValidationErrorResponse struct {
	Error struct {
		Code    int               `json:"code"`
		Message string            `json:"message"`
		Fields  map[string]string `json:"fields,omitempty"` // Поля с ошибками
	} `json:"error"`
}

// ErrorHandler глобальный middleware для обработки ошибок
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var statusCode int
			var message string
			var fieldErrors map[string]string // Словарь ошибок валидации

			var validationErrs validator.ValidationErrors // Объявляем переменную перед switch

			switch {
			case errors.Is(err, errs.ErrUserAlreadyExists):
				statusCode = http.StatusBadRequest
				message = "user already exists"
			case errors.Is(err, errs.ErrUserNotFound):
				statusCode = http.StatusNotFound
				message = "user not found"
			case errors.Is(err, errs.ErrEmailAlreadyUsed):
				statusCode = http.StatusBadRequest
				message = "email already used"
			case errors.Is(err, errs.ErrInvalidPwd):
				statusCode = http.StatusForbidden
				message = "invalid username or password"
			case errors.Is(err, errs.ErrInvalidTokenType):
				statusCode = http.StatusForbidden
				message = "invalid token type"
			case errors.Is(err, errs.ErrTokenExpired):
				statusCode = http.StatusForbidden
				message = "token invalid"
				fieldErrors = make(map[string]string)
				fieldErrors["refresh_token"] = errs.ErrTokenExpired.Error()
			case errors.Is(err, errs.ErrTokenInvalid):
				statusCode = http.StatusForbidden
				message = "token invalid"
				fieldErrors = make(map[string]string)
				fieldErrors["refresh_token"] = errs.ErrTokenInvalid.Error()
			case errors.Is(err, errs.ErrRefreshTokenRequired):
				statusCode = http.StatusBadRequest
				message = "token invalid"
				fieldErrors = make(map[string]string)
				fieldErrors["refresh_token"] = errs.ErrRefreshTokenRequired.Error()
			case errors.Is(err, errs.ErrValidateInRedis):
				statusCode = http.StatusForbidden
				message = "upss...something went wrong :("
			case errors.Is(err, errs.ErrTokenNotFound):
				statusCode = http.StatusBadRequest
				message = "token not found"
			case errors.Is(err, errs.ErrFailedToSave):
				statusCode = http.StatusBadRequest
				message = "try again later"
			case errors.Is(err, errs.ErrFailedToRefresh):
				statusCode = http.StatusBadRequest
				message = "try again later"
			case errors.As(err, &validationErrs): // Проверяем, является ли err ошибкой валидации
				statusCode = http.StatusBadRequest
				message = "Validation error"
				fieldErrors = make(map[string]string)
				for _, fieldErr := range validationErrs {
					fieldErrors[fieldErr.Field()] = validationErrorMessage(fieldErr)
				}

			default:
				statusCode = http.StatusInternalServerError
				message = "Internal server error"
			}

			// Логируем критические ошибки
			if statusCode == http.StatusInternalServerError {
				logger.Errorf("Unhandled server error: %v", err)
			}

			// Формируем JSON-ответ
			errorResponse := ValidationErrorResponse{}
			errorResponse.Error.Code = statusCode
			errorResponse.Error.Message = message
			if len(fieldErrors) > 0 {
				errorResponse.Error.Fields = fieldErrors
			}

			// Отправляем JSON-ответ
			c.JSON(statusCode, errorResponse)
		}
	}
}

// validationErrorMessage формирует читаемое сообщение ошибки
func validationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "min":
		return "must be at least " + fe.Param()
	case "max":
		return "must be at most " + fe.Param()
	case "gt":
		return "must be greater than " + fe.Param()
	case "email":
		return "must be a valid email address"
	default:
		return "is invalid"
	}
}
