package handler

import (
	"net/http"
	"strings"

	"chat-service/internal/usecase/session"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Middleware struct {
	sessionUsecase session.SessionUsecase
	logger         *logrus.Logger
}

func NewMiddleware(sessionUsecase session.SessionUsecase, logger *logrus.Logger) *Middleware {
	return &Middleware{
		sessionUsecase: sessionUsecase,
		logger:         logger,
	}
}

// AuthMiddleware проверяет JWT токен и устанавливает userID в контекст
func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			m.logger.Warn("authorization header is missing")
			SendError(c, "Authorization required", "No authorization header provided", http.StatusUnauthorized)
			c.Abort()
			return
		}

		// Ожидаем формат: "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			m.logger.Warn("invalid authorization header format")
			SendError(c, "Invalid authorization format", "Use 'Bearer <token>' format", http.StatusUnauthorized)
			c.Abort()
			return
		}

		// Валидируем сессию
		session, err := m.sessionUsecase.ValidateSession(c.Request.Context(), tokenString)
		if err != nil {
			m.logger.WithError(err).Warn("session validation failed")
			SendError(c, "Invalid session", "Session is invalid or expired", http.StatusUnauthorized)
			c.Abort()
			return
		}

		// Устанавливаем userID в контекст
		c.Set("userID", session.UserID)
		c.Set("session", session)
		c.Next()
	}
}

// CORSMiddleware добавляет CORS заголовки
func (m *Middleware) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

// LoggingMiddleware логирует каждый запрос
func (m *Middleware) LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// До обработки запроса
		m.logger.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"client": c.ClientIP(),
		}).Info("incoming request")

		// Обрабатываем запрос
		c.Next()

		// После обработки запроса
		m.logger.WithFields(logrus.Fields{
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
			"status":  c.Writer.Status(),
			"latency": c.Writer.Size(),
			"client":  c.ClientIP(),
		}).Info("request completed")
	}
}

// GetUserFromContext извлекает userID из контекста
func GetUserFromContext(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, &UnauthorizedErrorImpl{"user not authenticated"}
	}

	if id, ok := userID.(uuid.UUID); ok {
		return id, nil
	}

	return uuid.Nil, &UnauthorizedErrorImpl{"invalid user ID in context"}
}

type UnauthorizedErrorImpl struct {
	Message string
}

func (e *UnauthorizedErrorImpl) Error() string {
	return e.Message
}

func (e *UnauthorizedErrorImpl) Unauthorized() bool {
	return true
}
