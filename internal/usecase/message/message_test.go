package message

import (
	"context"
	"testing"
	"time"

	"chat-service/internal/entity"
	"chat-service/internal/usecase/mocks"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMessageUsecase_CreateMessage_Success(t *testing.T) {
	// Arrange
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel) // Отключаем логи в тестах

	messageRepo := &mocks.MessageRepoMock{}
	userRepo := &mocks.UserRepoMock{}

	testUserID := uuid.New()
	testContent := "Test message content"

	// Настраиваем моки
	userRepo.GetByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.User, error) {
		// Пользователь существует
		return &entity.User{ID: id}, nil
	}

	messageRepo.CreateFunc = func(ctx context.Context, message *entity.Message) error {
		// Успешное создание
		return nil
	}

	usecase := NewMessageUsecase(messageRepo, userRepo, logger)

	// Act
	message, err := usecase.CreateMessage(context.Background(), testUserID, testContent)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, testUserID, message.UserID)
	assert.Equal(t, testContent, message.Content)
	assert.NotEmpty(t, message.ID)
	assert.WithinDuration(t, time.Now(), message.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), message.UpdatedAt, time.Second)
}

func TestMessageUsecase_CreateMessage_UserNotFound(t *testing.T) {
	// Arrange
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	messageRepo := &mocks.MessageRepoMock{}
	userRepo := &mocks.UserRepoMock{}

	testUserID := uuid.New()
	testContent := "Test message content"

	// Настраиваем моки - пользователь не найден
	userRepo.GetByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.User, error) {
		return nil, &NotFoundError{"user not found"}
	}

	usecase := NewMessageUsecase(messageRepo, userRepo, logger)

	// Act
	message, err := usecase.CreateMessage(context.Background(), testUserID, testContent)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, message)
	assert.Contains(t, err.Error(), "user not found")
}

func TestMessageUsecase_CreateMessage_ValidationFailed(t *testing.T) {
	// Arrange
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	messageRepo := &mocks.MessageRepoMock{}
	userRepo := &mocks.UserRepoMock{}

	testUserID := uuid.New()
	invalidContent := "" // Пустой контент

	// Настраиваем моки - пользователь существует
	userRepo.GetByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.User, error) {
		return &entity.User{ID: id}, nil
	}

	usecase := NewMessageUsecase(messageRepo, userRepo, logger)

	// Act
	message, err := usecase.CreateMessage(context.Background(), testUserID, invalidContent)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, message)
	assert.Contains(t, err.Error(), "content is required")
}

func TestMessageUsecase_GetMessageByID_Success(t *testing.T) {
	// Arrange
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	messageRepo := &mocks.MessageRepoMock{}
	userRepo := &mocks.UserRepoMock{}

	testMessageID := uuid.New()
	expectedMessage := &entity.Message{
		ID:        testMessageID,
		UserID:    uuid.New(),
		Content:   "Test message",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Настраиваем моки
	messageRepo.GetByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.Message, error) {
		return expectedMessage, nil
	}

	usecase := NewMessageUsecase(messageRepo, userRepo, logger)

	// Act
	message, err := usecase.GetMessageByID(context.Background(), testMessageID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, expectedMessage, message)
}

func TestMessageUsecase_GetMessageByID_NotFound(t *testing.T) {
	// Arrange
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	messageRepo := &mocks.MessageRepoMock{}
	userRepo := &mocks.UserRepoMock{}

	testMessageID := uuid.New()

	// Настраиваем моки - сообщение не найдено
	messageRepo.GetByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.Message, error) {
		return nil, &NotFoundError{"message not found"}
	}

	usecase := NewMessageUsecase(messageRepo, userRepo, logger)

	// Act
	message, err := usecase.GetMessageByID(context.Background(), testMessageID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, message)
	assert.Contains(t, err.Error(), "message not found")
}

func TestMessageUsecase_GetMessagesByUser_Success(t *testing.T) {
	// Arrange
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	messageRepo := &mocks.MessageRepoMock{}
	userRepo := &mocks.UserRepoMock{}

	testUserID := uuid.New()
	messages := []*entity.Message{
		{
			ID:        uuid.New(),
			UserID:    testUserID,
			Content:   "Message 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			UserID:    testUserID,
			Content:   "Message 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Настраиваем моки
	userRepo.GetByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.User, error) {
		// Пользователь существует
		return &entity.User{ID: id}, nil
	}

	messageRepo.GetByUserIDFunc = func(ctx context.Context, userID uuid.UUID) ([]*entity.Message, error) {
		return messages, nil
	}

	usecase := NewMessageUsecase(messageRepo, userRepo, logger)

	// Act
	result, err := usecase.GetMessagesByUser(context.Background(), testUserID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, messages, result)
	assert.Len(t, result, 2)
}

func TestMessageUsecase_GetMessagesByUser_UserNotFound(t *testing.T) {
	// Arrange
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	messageRepo := &mocks.MessageRepoMock{}
	userRepo := &mocks.UserRepoMock{}

	testUserID := uuid.New()

	// Настраиваем моки - пользователь не найден
	userRepo.GetByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.User, error) {
		return nil, &NotFoundError{"user not found"}
	}

	usecase := NewMessageUsecase(messageRepo, userRepo, logger)

	// Act
	messages, err := usecase.GetMessagesByUser(context.Background(), testUserID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, messages)
	assert.Contains(t, err.Error(), "user not found")
}

func TestMessageUsecase_GetAllMessages_Success(t *testing.T) {
	// Arrange
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	messageRepo := &mocks.MessageRepoMock{}
	userRepo := &mocks.UserRepoMock{}

	messages := []*entity.Message{
		{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			Content:   "Public Message 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			Content:   "Public Message 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Настраиваем моки
	messageRepo.GetAllFunc = func(ctx context.Context) ([]*entity.Message, error) {
		return messages, nil
	}

	usecase := NewMessageUsecase(messageRepo, userRepo, logger)

	// Act
	result, err := usecase.GetAllMessages(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, messages, result)
	assert.Len(t, result, 2)
}

func TestMessageUsecase_DeleteMessage_Success(t *testing.T) {
	// Arrange
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	messageRepo := &mocks.MessageRepoMock{}
	userRepo := &mocks.UserRepoMock{}

	testMessageID := uuid.New()

	// Настраиваем моки
	messageRepo.DeleteFunc = func(ctx context.Context, id uuid.UUID) error {
		return nil // Успешное удаление
	}

	usecase := NewMessageUsecase(messageRepo, userRepo, logger)

	// Act
	err := usecase.DeleteMessage(context.Background(), testMessageID)

	// Assert
	assert.NoError(t, err)
}

// NotFoundError представляет ошибку, когда ресурс не найден.
type NotFoundError struct {
	Message string
}

// Error реализует интерфейс error.
func (e *NotFoundError) Error() string {
	return e.Message
}

// NotFound сигнализирует, что это ошибка "не найдено".
// Полезно для проверки типа в хендлерах или других местах.
func (e *NotFoundError) NotFound() bool {
	return true
}
