package mocks

import (
	"context"

	"chat-service/internal/entity"

	"github.com/google/uuid"
)

type MessageRepoMock struct {
	CreateFunc      func(ctx context.Context, message *entity.Message) error
	GetByIDFunc     func(ctx context.Context, id uuid.UUID) (*entity.Message, error)
	GetByUserIDFunc func(ctx context.Context, userID uuid.UUID) ([]*entity.Message, error)
	GetAllFunc      func(ctx context.Context) ([]*entity.Message, error)
	DeleteFunc      func(ctx context.Context, id uuid.UUID) error
}

func (m *MessageRepoMock) Create(ctx context.Context, message *entity.Message) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, message)
	}
	return nil
}

func (m *MessageRepoMock) GetByID(ctx context.Context, id uuid.UUID) (*entity.Message, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MessageRepoMock) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Message, error) {
	if m.GetByUserIDFunc != nil {
		return m.GetByUserIDFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MessageRepoMock) GetAll(ctx context.Context) ([]*entity.Message, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return nil, nil
}

func (m *MessageRepoMock) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}
