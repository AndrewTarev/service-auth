package http

import (
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
)

type SuccessResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Message string `json:"error"`
}

func newErrorResponse(c *gin.Context, statusCode int, logMessage, userMessage string) {
	logger.Error(logMessage)

	errJSON := ErrorResponse{Message: userMessage}
	c.AbortWithStatusJSON(statusCode, errJSON)
}
