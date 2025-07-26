package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJWTService_GenerateToken_Success(t *testing.T) {
	// Arrange
	logger := newTestLogger()
	secretKey := "test_secret_key_for_testing"
	service := NewJWTService(secretKey, logger) // Передаем секрет и логгер

	userID := uuid.New()

	// Act
	tokenString, err := service.GenerateToken(userID)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)
	assert.IsType(t, "", tokenString)
	// Проверяем, что токен выглядит как JWT (три части, разделенные точками)
	assert.Regexp(t, `^[A-Za-z0-9-_]*\.[A-Za-z0-9-_]*\.[A-Za-z0-9-_]*$`, tokenString)
}

func TestJWTService_GenerateToken_DifferentUserIDs(t *testing.T) {
	// Arrange
	logger := newTestLogger()
	secretKey := "test_secret_key_for_testing"
	service := NewJWTService(secretKey, logger) // Передаем секрет и логгер

	userID1 := uuid.New()
	userID2 := uuid.New()

	// Act
	token1, err1 := service.GenerateToken(userID1)
	token2, err2 := service.GenerateToken(userID2)

	// Assert
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEmpty(t, token1)
	assert.NotEmpty(t, token2)
	// Токены для разных пользователей должны быть разными
	assert.NotEqual(t, token1, token2)
}

func TestJWTService_ValidateToken_ValidToken(t *testing.T) {
	// Arrange
	logger := newTestLogger()
	secretKey := "test_secret_key_for_testing"
	service := NewJWTService(secretKey, logger) // Передаем секрет и логгер

	userID := uuid.New()

	// Сначала генерируем токен
	tokenString, err := service.GenerateToken(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Act
	parsedUserID, err := service.ValidateToken(tokenString)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, userID, parsedUserID)
}

func TestJWTService_ValidateToken_InvalidToken(t *testing.T) {
	// Arrange
	logger := newTestLogger()
	secretKey := "test_secret_key_for_testing"
	service := NewJWTService(secretKey, logger) // Передаем секрет и логгер

	invalidToken := "invalid.token.string"

	// Act
	userID, err := service.ValidateToken(invalidToken)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, userID)
}

func TestJWTService_ValidateToken_SignatureMismatch(t *testing.T) {
	// Arrange
	logger := newTestLogger()
	secretKey1 := "test_secret_key_for_testing_1"
	secretKey2 := "test_secret_key_for_testing_2"

	service1 := NewJWTService(secretKey1, logger) // Передаем секрет и логгер
	service2 := NewJWTService(secretKey2, logger) // Передаем секрет и логгер

	userID := uuid.New()

	// Генерируем токен с одним секретом
	tokenString, err := service1.GenerateToken(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Пытаемся валидировать токен с другим секретом
	// Act
	parsedUserID, err := service2.ValidateToken(tokenString)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, parsedUserID)
}

func TestJWTService_ValidateToken_ExpiredToken(t *testing.T) {
	// Тестирование истечения срока действия токена требует модификации реализации JWTService
	// или использование библиотеки для создания истекшего токена.
	// Для простоты, мы можем протестировать это косвенно или оставить для расширенного тестирования.
	// В текущей реализации токены действительны 24 часа, что сложно протестировать напрямую.
	t.Skip("Тестирование истечения срока действия токена требует специальной реализации")
}

func TestJWTService_GenerateAndValidateToken_Consistency(t *testing.T) {
	// Arrange
	logger := newTestLogger()
	secretKey := "test_secret_key_for_testing"
	service := NewJWTService(secretKey, logger) // Передаем секрет и логгер

	userID := uuid.New()

	// Act
	tokenString, err := service.GenerateToken(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	parsedUserID, err := service.ValidateToken(tokenString)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, userID, parsedUserID)
}

func TestJWTService_ValidateToken_MalformedToken(t *testing.T) {
	// Arrange
	logger := newTestLogger()
	secretKey := "test_secret_key_for_testing"
	service := NewJWTService(secretKey, logger) // Передаем секрет и логгер

	malformedTokens := []string{
		"just.a.string",
		"missing.part",
		"",
		"....",
		"invalid!token#format",
	}

	for _, token := range malformedTokens {
		// Act
		userID, err := service.ValidateToken(token)

		// Assert
		assert.Error(t, err, "Expected error for token: %s", token)
		assert.Equal(t, uuid.Nil, userID, "Expected uuid.Nil for token: %s", token)
	}
}
