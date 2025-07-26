package handler

import (
	"net/http"

	"chat-service/internal/entity"
	"chat-service/internal/usecase/message"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type MessageHandler struct {
	messageUsecase message.MessageUsecase
	logger         *logrus.Logger
}

func NewMessageHandler(
	messageUsecase message.MessageUsecase,
	logger *logrus.Logger,
) *MessageHandler {
	return &MessageHandler{
		messageUsecase: messageUsecase,
		logger:         logger,
	}
}

// CreateMessageRequest структура для создания сообщения
// swagger:model CreateMessageRequest
type CreateMessageRequest struct {
	// Текст сообщения
	// required: true
	// min length: 1
	// max length: 1000
	Content string `json:"content" binding:"required,min=1,max=1000"`
}

// MessageResponse структура ответа с сообщением
// swagger:model MessageResponse
type MessageResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    *entity.Message `json:"data"`
}

// MessagesResponse структура ответа с массивом сообщений
// swagger:model MessagesResponse
type MessagesResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    []*entity.Message `json:"data"`
}

// CreateMessage создает новое сообщение
// @Summary Создание нового сообщения
// @Description Создает новое сообщение от авторизованного пользователя
// @Tags messages
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param message body CreateMessageRequest true "Текст сообщения"
// @Success 201 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /messages [post]
func (h *MessageHandler) CreateMessage(c *gin.Context) {
	userID, err := GetUserFromContext(c)
	if err != nil {
		h.logger.WithError(err).Warn("failed to get user from context")
		HandleError(c, err, h.logger)
		return
	}

	var req CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Warn("invalid create message request body")
		SendError(c, "Invalid request", err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id": userID,
		"content": req.Content[:min(50, len(req.Content))] + "...",
	}).Info("creating new message")

	message, err := h.messageUsecase.CreateMessage(c.Request.Context(), userID, req.Content)
	if err != nil {
		h.logger.WithError(err).Error("failed to create message")
		HandleError(c, err, h.logger)
		return
	}

	h.logger.WithField("message_id", message.ID).Info("message created successfully")
	SendSuccess(c, message, "Message created successfully", http.StatusCreated)
}

// GetMessageByID возвращает сообщение по ID
// @Summary Получение сообщения по ID
// @Description Возвращает конкретное сообщение по его идентификатору
// @Tags messages
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param id path string true "ID сообщения" Format(uuid)
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /messages/{id} [get]
func (h *MessageHandler) GetMessageByID(c *gin.Context) {
	messageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.WithError(err).Warn("invalid message ID format")
		SendError(c, "Invalid message ID", "Message ID must be a valid UUID", http.StatusBadRequest)
		return
	}

	h.logger.WithField("message_id", messageID).Debug("fetching message by ID")

	message, err := h.messageUsecase.GetMessageByID(c.Request.Context(), messageID)
	if err != nil {
		h.logger.WithError(err).Error("failed to fetch message by ID")
		HandleError(c, err, h.logger)
		return
	}

	h.logger.WithField("message_id", messageID).Debug("message fetched successfully")
	SendSuccess(c, message, "Message retrieved successfully", http.StatusOK)
}

// GetMessagesByUser возвращает все сообщения пользователя
// @Summary Получение всех сообщений пользователя
// @Description Возвращает все сообщения авторизованного пользователя
// @Tags messages
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {object} MessagesResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /messages/my [get]
func (h *MessageHandler) GetMessagesByUser(c *gin.Context) {
	userID, err := GetUserFromContext(c)
	if err != nil {
		h.logger.WithError(err).Warn("failed to get user from context")
		HandleError(c, err, h.logger)
		return
	}

	h.logger.WithField("user_id", userID).Debug("fetching messages for user")

	messages, err := h.messageUsecase.GetMessagesByUser(c.Request.Context(), userID)
	if err != nil {
		h.logger.WithError(err).Error("failed to fetch user messages")
		HandleError(c, err, h.logger)
		return
	}

	h.logger.WithField("user_id", userID).Debugf("fetched %d messages for user", len(messages))
	SendSuccess(c, messages, "Messages retrieved successfully", http.StatusOK)
}

// GetAllMessages возвращает все сообщения
// @Summary Получение всех сообщений
// @Description Возвращает все сообщения в системе (публичный доступ)
// @Tags messages
// @Accept  json
// @Produce  json
// @Success 200 {object} MessagesResponse
// @Failure 500 {object} ErrorResponse
// @Router /messages [get]
func (h *MessageHandler) GetAllMessages(c *gin.Context) {
	h.logger.Debug("fetching all messages")

	messages, err := h.messageUsecase.GetAllMessages(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("failed to fetch all messages")
		HandleError(c, err, h.logger)
		return
	}

	h.logger.Debugf("fetched %d messages total", len(messages))
	SendSuccess(c, messages, "Messages retrieved successfully", http.StatusOK)
}

// DeleteMessage удаляет сообщение
// @Summary Удаление сообщения
// @Description Удаляет сообщение авторизованного пользователя
// @Tags messages
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param id path string true "ID сообщения" Format(uuid)
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /messages/{id} [delete]
func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	userID, err := GetUserFromContext(c)
	if err != nil {
		h.logger.WithError(err).Warn("failed to get user from context")
		HandleError(c, err, h.logger)
		return
	}

	messageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.WithError(err).Warn("invalid message ID format")
		SendError(c, "Invalid message ID", "Message ID must be a valid UUID", http.StatusBadRequest)
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"message_id": messageID,
	}).Warn("message deletion requested")

	// TODO: Проверить права доступа (владелец сообщения или админ)
	// Пока что разрешаем владельцу удалять свои сообщения

	// Получаем сообщение для проверки владельца
	message, err := h.messageUsecase.GetMessageByID(c.Request.Context(), messageID)
	if err != nil {
		h.logger.WithError(err).Error("failed to fetch message for deletion check")
		HandleError(c, err, h.logger)
		return
	}

	if message.UserID != userID {
		h.logger.WithFields(logrus.Fields{
			"user_id":    userID,
			"message_id": messageID,
			"owner_id":   message.UserID,
		}).Warn("user trying to delete another user's message")
		SendError(c, "Forbidden", "You can only delete your own messages", http.StatusForbidden)
		return
	}

	err = h.messageUsecase.DeleteMessage(c.Request.Context(), messageID)
	if err != nil {
		h.logger.WithError(err).Error("failed to delete message")
		HandleError(c, err, h.logger)
		return
	}

	h.logger.WithField("message_id", messageID).Info("message deleted successfully")
	SendSuccess(c, nil, "Message deleted successfully", http.StatusOK)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
