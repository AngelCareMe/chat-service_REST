package user

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

func TestUserUsecase_Register_Success(t *testing.T) {
	// Arrange
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel) // Отключаем логи в тестах

	userRepo := &mocks.UserRepoMock{}
	sessionRepo := &mocks.SessionRepoMock{}
	hashService := &mocks.HashServiceMock{}
	jwtService := &mocks.JWTServiceMock{}

	// Настраиваем моки
	userRepo.GetByEmailFunc = func(ctx context.Context, email string) (*entity.User, error) {
		return nil, &NotFoundError{"user not found"} // Пользователь не существует
	}

	userRepo.CreateFunc = func(ctx context.Context, user *entity.User) error {
		return nil // Успешное создание
	}

	hashService.HashPasswordFunc = func(password string) (string, error) {
		return "hashed_password", nil
	}

	usecase := NewUserUsecase(userRepo, sessionRepo, hashService, jwtService, logger)

	// Act
	user, err := usecase.Register(context.Background(), "testuser", "test@example.com", "password123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Empty(t, user.Password) // Пароль должен быть очищен
}

func TestUserUsecase_Register_UserExists(t *testing.T) {
	// Arrange
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	userRepo := &mocks.UserRepoMock{}
	sessionRepo := &mocks.SessionRepoMock{}
	hashService := &mocks.HashServiceMock{}
	jwtService := &mocks.JWTServiceMock{}

	// Настраиваем моки - пользователь уже существует
	existingUser := &entity.User{
		ID:       uuid.New(),
		Username: "existinguser",
		Email:    "test@example.com",
		Password: "hashed_password",
	}

	userRepo.GetByEmailFunc = func(ctx context.Context, email string) (*entity.User, error) {
		return existingUser, nil // Пользователь существует
	}

	usecase := NewUserUsecase(userRepo, sessionRepo, hashService, jwtService, logger)

	// Act
	user, err := usecase.Register(context.Background(), "testuser", "test@example.com", "password123")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "user with this email already exists")
}

func TestUserUsecase_Login_Success(t *testing.T) {
	// Arrange
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	userRepo := &mocks.UserRepoMock{}
	sessionRepo := &mocks.SessionRepoMock{}
	hashService := &mocks.HashServiceMock{}
	jwtService := &mocks.JWTServiceMock{}

	// Настраиваем моки
	testUser := &entity.User{
		ID:        uuid.New(),
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashed_password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	userRepo.GetByEmailFunc = func(ctx context.Context, email string) (*entity.User, error) {
		return testUser, nil
	}

	hashService.CheckPasswordHashFunc = func(password, hash string) bool {
		return true // Правильный пароль
	}

	usecase := NewUserUsecase(userRepo, sessionRepo, hashService, jwtService, logger)

	// Act
	user, err := usecase.Login(context.Background(), "test@example.com", "password123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Empty(t, user.Password) // Пароль должен быть очищен
}

func TestUserUsecase_Login_InvalidCredentials(t *testing.T) {
	// Arrange
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	userRepo := &mocks.UserRepoMock{}
	sessionRepo := &mocks.SessionRepoMock{}
	hashService := &mocks.HashServiceMock{}
	jwtService := &mocks.JWTServiceMock{}

	// Настраиваем моки - пользователь не найден
	userRepo.GetByEmailFunc = func(ctx context.Context, email string) (*entity.User, error) {
		return nil, &NotFoundError{"user not found"}
	}

	usecase := NewUserUsecase(userRepo, sessionRepo, hashService, jwtService, logger)

	// Act
	user, err := usecase.Login(context.Background(), "test@example.com", "wrongpassword")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestUserUsecase_GetProfile_Success(t *testing.T) {
	// Arrange
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	userRepo := &mocks.UserRepoMock{}
	sessionRepo := &mocks.SessionRepoMock{}
	hashService := &mocks.HashServiceMock{}
	jwtService := &mocks.JWTServiceMock{}

	testUserID := uuid.New()
	testUser := &entity.User{
		ID:        testUserID,
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashed_password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	userRepo.GetByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.User, error) {
		return testUser, nil
	}

	usecase := NewUserUsecase(userRepo, sessionRepo, hashService, jwtService, logger)

	// Act
	user, err := usecase.GetProfile(context.Background(), testUserID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, testUserID, user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Empty(t, user.Password) // Пароль должен быть очищен
}

func TestUserUsecase_UpdateProfile_Success(t *testing.T) {
	// Arrange
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	userRepo := &mocks.UserRepoMock{}
	sessionRepo := &mocks.SessionRepoMock{}
	hashService := &mocks.HashServiceMock{}
	jwtService := &mocks.JWTServiceMock{}

	testUser := &entity.User{
		ID:        uuid.New(),
		Username:  "updateduser",
		Email:     "updated@example.com",
		Password:  "", // Пароль не обязателен при обновлении
		UpdatedAt: time.Now(),
	}

	userRepo.UpdateFunc = func(ctx context.Context, user *entity.User) error {
		return nil // Успешное обновление
	}

	usecase := NewUserUsecase(userRepo, sessionRepo, hashService, jwtService, logger)

	// Act
	err := usecase.UpdateProfile(context.Background(), testUser)

	// Assert
	assert.NoError(t, err)
}

func TestUserUsecase_DeleteUser_Success(t *testing.T) {
	// Arrange
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	userRepo := &mocks.UserRepoMock{}
	sessionRepo := &mocks.SessionRepoMock{}
	hashService := &mocks.HashServiceMock{}
	jwtService := &mocks.JWTServiceMock{}

	testUserID := uuid.New()

	sessionRepo.GetByUserIDFunc = func(ctx context.Context, userID uuid.UUID) (*entity.Session, error) {
		return nil, &NotFoundError{"session not found"} // Сессия не найдена
	}

	userRepo.DeleteFunc = func(ctx context.Context, id uuid.UUID) error {
		return nil // Успешное удаление
	}

	usecase := NewUserUsecase(userRepo, sessionRepo, hashService, jwtService, logger)

	// Act
	err := usecase.DeleteUser(context.Background(), testUserID)

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
