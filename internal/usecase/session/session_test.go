package session

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

func TestSessionUsecase_CreateSession_Success(t *testing.T) {
	// Arrange
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel) // Отключаем логи в тестах

	sessionRepo := &mocks.SessionRepoMock{}
	jwtService := &mocks.JWTServiceMock{}

	testUserID := uuid.New()
	testToken := "generated_jwt_token"

	// Настраиваем моки
	jwtService.GenerateTokenFunc = func(userID uuid.UUID) (string, error) {
		return testToken, nil // Успешная генерация токена
	}

	sessionRepo.CreateFunc = func(ctx context.Context, session *entity.Session) error {
		return nil // Успешное создание сессии
	}

	usecase := NewSessionUsecase(sessionRepo, jwtService, logger)

	// Act
	session, err := usecase.CreateSession(context.Background(), testUserID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, testUserID, session.UserID)
	assert.Equal(t, testToken, session.Token)
	assert.NotEmpty(t, session.ID)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), session.ExpiresAt, time.Minute) // Проверяем, что срок ~24 часа
	assert.WithinDuration(t, time.Now(), session.CreatedAt, time.Second)
}

func TestSessionUsecase_CreateSession_JWTError(t *testing.T) {
	// Arrange
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	sessionRepo := &mocks.SessionRepoMock{}
	jwtService := &mocks.JWTServiceMock{}

	testUserID := uuid.New()

	// Настраиваем моки - ошибка генерации токена
	jwtService.GenerateTokenFunc = func(userID uuid.UUID) (string, error) {
		return "", &BusinessError{"failed to generate token"}
	}

	usecase := NewSessionUsecase(sessionRepo, jwtService, logger)

	// Act
	session, err := usecase.CreateSession(context.Background(), testUserID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, session)
	assert.Contains(t, err.Error(), "failed to generate token")
}

func TestSessionUsecase_ValidateSession_Success(t *testing.T) {
	// Arrange
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	sessionRepo := &mocks.SessionRepoMock{}
	jwtService := &mocks.JWTServiceMock{}

	testToken := "valid_token"
	testUserID := uuid.New()
	testSession := &entity.Session{
		ID:        uuid.New(),
		UserID:    testUserID,
		Token:     testToken,
		ExpiresAt: time.Now().Add(time.Hour), // Не истек
		CreatedAt: time.Now(),
	}

	// Настраиваем моки
	sessionRepo.GetByTokenFunc = func(ctx context.Context, token string) (*entity.Session, error) {
		return testSession, nil // Сессия найдена
	}

	jwtService.ValidateTokenFunc = func(token string) (uuid.UUID, error) {
		return testUserID, nil // Токен валиден
	}

	usecase := NewSessionUsecase(sessionRepo, jwtService, logger)

	// Act
	session, err := usecase.ValidateSession(context.Background(), testToken)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, testSession, session)
}

func TestSessionUsecase_ValidateSession_SessionNotFound(t *testing.T) {
	// Arrange
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	sessionRepo := &mocks.SessionRepoMock{}
	jwtService := &mocks.JWTServiceMock{}

	testToken := "invalid_token"

	// Настраиваем моки - сессия не найдена
	sessionRepo.GetByTokenFunc = func(ctx context.Context, token string) (*entity.Session, error) {
		return nil, &NotFoundError{"session not found"}
	}

	usecase := NewSessionUsecase(sessionRepo, jwtService, logger)

	// Act
	session, err := usecase.ValidateSession(context.Background(), testToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, session)
	assert.Contains(t, err.Error(), "invalid session")
}

func TestSessionUsecase_ValidateSession_SessionExpired(t *testing.T) {
	// Arrange
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	sessionRepo := &mocks.SessionRepoMock{}
	jwtService := &mocks.JWTServiceMock{}

	testToken := "expired_token"
	testSession := &entity.Session{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Token:     testToken,
		ExpiresAt: time.Now().Add(-time.Hour), // Истекла
		CreatedAt: time.Now().Add(-2 * time.Hour),
	}

	// Настраиваем моки - сессия найдена, но истекла
	sessionRepo.GetByTokenFunc = func(ctx context.Context, token string) (*entity.Session, error) {
		return testSession, nil
	}

	// Ожидаем, что сессия будет удалена
	sessionRepo.DeleteByTokenFunc = func(ctx context.Context, token string) error {
		return nil // Успешное удаление
	}

	usecase := NewSessionUsecase(sessionRepo, jwtService, logger)

	// Act
	session, err := usecase.ValidateSession(context.Background(), testToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, session)
	assert.Contains(t, err.Error(), "session expired")
}

func TestSessionUsecase_ValidateSession_InvalidJWT(t *testing.T) {
	// Arrange
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	sessionRepo := &mocks.SessionRepoMock{}
	jwtService := &mocks.JWTServiceMock{}

	testToken := "token_with_invalid_jwt"
	testSession := &entity.Session{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Token:     testToken,
		ExpiresAt: time.Now().Add(time.Hour), // Не истек
		CreatedAt: time.Now(),
	}

	// Настраиваем моки - сессия найдена, но JWT невалиден
	sessionRepo.GetByTokenFunc = func(ctx context.Context, token string) (*entity.Session, error) {
		return testSession, nil
	}

	jwtService.ValidateTokenFunc = func(token string) (uuid.UUID, error) {
		return uuid.Nil, &BusinessError{"invalid token"}
	}

	usecase := NewSessionUsecase(sessionRepo, jwtService, logger)

	// Act
	session, err := usecase.ValidateSession(context.Background(), testToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, session)
	assert.Contains(t, err.Error(), "invalid token")
}

func TestSessionUsecase_DeleteSession_Success(t *testing.T) {
	// Arrange
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	sessionRepo := &mocks.SessionRepoMock{}
	jwtService := &mocks.JWTServiceMock{}

	testToken := "token_to_delete"

	// Настраиваем моки
	sessionRepo.DeleteByTokenFunc = func(ctx context.Context, token string) error {
		return nil // Успешное удаление
	}

	usecase := NewSessionUsecase(sessionRepo, jwtService, logger)

	// Act
	err := usecase.DeleteSession(context.Background(), testToken)

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
