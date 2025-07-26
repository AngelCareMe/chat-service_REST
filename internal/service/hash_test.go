package service

import (
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func newTestLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(io.Discard)       // Правильный способ отключить вывод
	logger.SetLevel(logrus.FatalLevel) // Минимальный уровень логирования
	return logger
}

func TestHashService_HashPassword_Success(t *testing.T) {
	// Arrange
	logger := newTestLogger()
	service := NewHashService(logger) // Передаем логгер

	password := "test_password_123"

	// Act
	hashedPassword, err := service.HashPassword(password)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, password, hashedPassword)
	// Проверяем, что хэш выглядит как bcrypt хэш (начинается с $2a$, $2b$ или $2y$)
	assert.Regexp(t, `^\$2[aby]\$`, hashedPassword)
}

func TestHashService_HashPassword_EmptyPassword(t *testing.T) {
	// Arrange
	logger := newTestLogger()
	service := NewHashService(logger) // Передаем логгер

	password := ""

	// Act
	hashedPassword, err := service.HashPassword(password)

	// Assert
	assert.NoError(t, err) // bcrypt может хэшировать пустую строку
	assert.NotEmpty(t, hashedPassword)
}

func TestHashService_CheckPasswordHash_ValidPassword(t *testing.T) {
	// Arrange
	logger := newTestLogger()
	service := NewHashService(logger) // Передаем логгер

	password := "correct_password_123"

	// Сначала хэшируем пароль
	hashedPassword, err := service.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	// Act
	isValid := service.CheckPasswordHash(password, hashedPassword)

	// Assert
	assert.True(t, isValid)
}

func TestHashService_CheckPasswordHash_InvalidPassword(t *testing.T) {
	// Arrange
	logger := newTestLogger()
	service := NewHashService(logger) // Передаем логгер

	password := "correct_password_123"
	wrongPassword := "wrong_password_123"

	// Сначала хэшируем пароль
	hashedPassword, err := service.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	// Act
	isValid := service.CheckPasswordHash(wrongPassword, hashedPassword)

	// Assert
	assert.False(t, isValid)
}

func TestHashService_CheckPasswordHash_InvalidHash(t *testing.T) {
	// Arrange
	logger := newTestLogger()
	service := NewHashService(logger) // Передаем логгер

	password := "any_password"
	invalidHash := "invalid_hash_format"

	// Act
	isValid := service.CheckPasswordHash(password, invalidHash)

	// Assert
	assert.False(t, isValid)
}

func TestHashService_CheckPasswordHash_EmptyInputs(t *testing.T) {
	// Arrange
	logger := newTestLogger()
	service := NewHashService(logger) // Передаем логгер

	// Act & Assert - пустой пароль и пустой хэш
	isValid := service.CheckPasswordHash("", "")
	assert.False(t, isValid)

	// Пустой пароль с валидным хэшем
	validHash, _ := service.HashPassword("test")
	isValid = service.CheckPasswordHash("", validHash)
	assert.False(t, isValid)

	// Валидный пароль с пустым хэшем
	isValid = service.CheckPasswordHash("test", "")
	assert.False(t, isValid)
}
