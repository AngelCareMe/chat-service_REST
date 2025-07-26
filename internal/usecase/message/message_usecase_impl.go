package message

import (
	"chat-service/internal/entity"
	"chat-service/internal/usecase"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type messageUsecase struct {
	messageRepo usecase.MessageRepository
	userRepo    usecase.UserRepository
	logger      *logrus.Logger
}

func NewMessageUsecase(messageRepo usecase.MessageRepository, userRepo usecase.UserRepository, logger *logrus.Logger) MessageUsecase {
	return &messageUsecase{
		messageRepo: messageRepo,
		userRepo:    userRepo,
		logger:      logger,
	}
}

func (m *messageUsecase) CreateMessage(ctx context.Context, userID uuid.UUID, content string) (*entity.Message, error) {
	m.logger.WithFields(logrus.Fields{
		"user_id": userID,
		"content": content[:min(50, len(content))],
	}).Info("creating new message")

	// Проверяем существование пользователя
	m.logger.WithField("user_id", userID).Debug("checking user existence")
	_, err := m.userRepo.GetByID(ctx, userID)
	if err != nil {
		m.logger.WithError(err).WithField("user_id", userID).Warn("user not found")
		return nil, &BusinessError{"user not found"}
	}

	message := &entity.Message{
		ID:        uuid.New(),
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := message.Validate(); err != nil {
		m.logger.WithError(err).Warn("message validation failed")
		return nil, err
	}

	m.logger.WithField("message_id", message.ID).Debug("saving message to repository")
	if err := m.messageRepo.Create(ctx, message); err != nil {
		m.logger.WithError(err).WithField("message_id", message.ID).Error("failed to create message")
		return nil, err
	}

	m.logger.WithField("message_id", message.ID).Info("message created successfully")
	return message, nil
}

func (m *messageUsecase) GetMessageByID(ctx context.Context, messageID uuid.UUID) (*entity.Message, error) {
	m.logger.WithField("message_id", messageID).Debug("fetching message by ID")

	message, err := m.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		m.logger.WithError(err).WithField("message_id", messageID).Error("failed to fetch message")
		return nil, err
	}

	m.logger.WithField("message_id", messageID).Debug("message fetched successfully")
	return message, nil
}

func (m *messageUsecase) GetMessagesByUser(ctx context.Context, userID uuid.UUID) ([]*entity.Message, error) {
	m.logger.WithField("user_id", userID).Debug("fetching messages by user")

	// Проверяем существование пользователя
	m.logger.WithField("user_id", userID).Debug("checking user existence")
	_, err := m.userRepo.GetByID(ctx, userID)
	if err != nil {
		m.logger.WithError(err).WithField("user_id", userID).Warn("user not found")
		return nil, &BusinessError{"user not found"}
	}

	messages, err := m.messageRepo.GetByUserID(ctx, userID)
	if err != nil {
		m.logger.WithError(err).WithField("user_id", userID).Error("failed to fetch user messages")
		return nil, err
	}

	m.logger.WithField("user_id", userID).Debugf("fetched %d messages for user", len(messages))
	return messages, nil
}

func (m *messageUsecase) GetAllMessages(ctx context.Context) ([]*entity.Message, error) {
	m.logger.Debug("fetching all messages")

	messages, err := m.messageRepo.GetAll(ctx)
	if err != nil {
		m.logger.WithError(err).Error("failed to fetch all messages")
		return nil, err
	}

	m.logger.Debugf("fetched %d messages total", len(messages))
	return messages, nil
}

func (m *messageUsecase) DeleteMessage(ctx context.Context, messageID uuid.UUID) error {
	m.logger.WithField("message_id", messageID).Warn("deleting message")

	err := m.messageRepo.Delete(ctx, messageID)
	if err != nil {
		m.logger.WithError(err).WithField("message_id", messageID).Error("failed to delete message")
		return err
	}

	m.logger.WithField("message_id", messageID).Info("message deleted successfully")
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type BusinessError struct {
	Message string
}

func (e *BusinessError) Error() string {
	return e.Message
}

func (e *BusinessError) ValidationError() bool {
	return true
}
