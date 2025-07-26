package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewSuccessResponse(data interface{}, message string) *SuccessResponse {
	return &SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func NewErrorResponse(message, errorStr string) *ErrorResponse {
	return &ErrorResponse{
		Success: false,
		Message: message,
		Error:   errorStr,
	}
}

func SendSuccess(c *gin.Context, data interface{}, message string, statusCode int) {
	c.JSON(statusCode, NewSuccessResponse(data, message))
}

func SendError(c *gin.Context, message, errorStr string, statusCode int) {
	c.JSON(statusCode, NewErrorResponse(message, errorStr))
}

// Обработчик ошибок по типам
func HandleError(c *gin.Context, err error, logger *logrus.Logger) {
	logger.WithError(err).Error("handler error occurred")

	// Проверяем тип ошибки через type assertion
	switch e := err.(type) {
	case ValidationError:
		SendError(c, "Validation failed", e.Error(), http.StatusBadRequest)
	case NotFoundError:
		SendError(c, "Resource not found", e.Error(), http.StatusNotFound)
	case UnauthorizedError:
		SendError(c, "Unauthorized", e.Error(), http.StatusUnauthorized)
	default:
		SendError(c, "Internal server error", "Something went wrong", http.StatusInternalServerError)
	}
}

// Интерфейсы для типизации ошибок
type ValidationError interface {
	ValidationError() bool
	Error() string
}

type NotFoundError interface {
	NotFound() bool
	Error() string
}

type UnauthorizedError interface {
	Unauthorized() bool
	Error() string
}
