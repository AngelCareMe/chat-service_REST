package mocks

import (
	"github.com/google/uuid"
)

type JWTServiceMock struct {
	GenerateTokenFunc func(userID uuid.UUID) (string, error)
	ValidateTokenFunc func(token string) (uuid.UUID, error)
}

func (m *JWTServiceMock) GenerateToken(userID uuid.UUID) (string, error) {
	if m.GenerateTokenFunc != nil {
		return m.GenerateTokenFunc(userID)
	}
	return "test_token", nil
}

func (m *JWTServiceMock) ValidateToken(token string) (uuid.UUID, error) {
	if m.ValidateTokenFunc != nil {
		return m.ValidateTokenFunc(token)
	}
	return uuid.New(), nil
}
