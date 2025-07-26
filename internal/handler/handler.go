package handler

import (
	"net/http"

	"chat-service/internal/usecase/message"
	"chat-service/internal/usecase/session"
	"chat-service/internal/usecase/user"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Документация будет сгенерирована swag
	_ "chat-service/internal/docs"
)

type Handler struct {
	router         *gin.Engine
	userHandler    *UserHandler
	messageHandler *MessageHandler
	middleware     *Middleware
	logger         *logrus.Logger
}

func NewHandler(
	userUsecase user.UserUsecase,
	messageUsecase message.MessageUsecase,
	sessionUsecase session.SessionUsecase,
	logger *logrus.Logger,
) *Handler {
	// Устанавливаем режим Gin
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Middleware
	middleware := NewMiddleware(sessionUsecase, logger)

	// Handlers
	userHandler := NewUserHandler(userUsecase, sessionUsecase, logger)
	messageHandler := NewMessageHandler(messageUsecase, logger)

	handler := &Handler{
		router:         router,
		userHandler:    userHandler,
		messageHandler: messageHandler,
		middleware:     middleware,
		logger:         logger,
	}

	handler.setupRoutes()
	return handler
}

func (h *Handler) setupRoutes() {
	// Global middleware
	h.router.Use(h.middleware.LoggingMiddleware())
	h.router.Use(h.middleware.CORSMiddleware())
	h.router.Use(gin.Recovery())

	// Health check
	h.router.GET("/health", func(c *gin.Context) {
		SendSuccess(c, gin.H{"status": "ok"}, "Service is running", http.StatusOK)
	})

	// Swagger documentation
	h.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	public := h.router.Group("/api/v1")
	{
		public.POST("/register", h.userHandler.Register)
		public.POST("/login", h.userHandler.Login)
		public.GET("/messages", h.messageHandler.GetAllMessages)
	}

	// Protected routes
	protected := h.router.Group("/api/v1")
	protected.Use(h.middleware.AuthMiddleware())
	{
		protected.GET("/profile", h.userHandler.GetProfile)
		protected.PUT("/profile", h.userHandler.UpdateProfile)
		protected.POST("/logout", h.userHandler.Logout)
		protected.DELETE("/profile", h.userHandler.DeleteUser)
		protected.POST("/messages", h.messageHandler.CreateMessage)
		protected.GET("/messages/my", h.messageHandler.GetMessagesByUser)
		protected.GET("/messages/:id", h.messageHandler.GetMessageByID)
		protected.DELETE("/messages/:id", h.messageHandler.DeleteMessage)
	}

	h.logger.Info("routes configured successfully")
}

func (h *Handler) GetRouter() *gin.Engine {
	return h.router
}

// Close освобождает ресурсы handler'а
func (h *Handler) Close() {
	h.logger.Info("closing handler resources")
	// Здесь можно закрыть дополнительные ресурсы, если появятся
}
