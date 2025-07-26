package mocks

import (
	"context"

	"chat-service/internal/entity"

	"github.com/google/uuid"
)

type SessionRepoMock struct {
	CreateFunc        func(ctx context.Context, session *entity.Session) error
	GetByTokenFunc    func(ctx context.Context, token string) (*entity.Session, error)
	GetByUserIDFunc   func(ctx context.Context, userID uuid.UUID) (*entity.Session, error)
	DeleteFunc        func(ctx context.Context, id uuid.UUID) error
	DeleteByTokenFunc func(ctx context.Context, token string) error
}

func (m *SessionRepoMock) Create(ctx context.Context, session *entity.Session) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, session)
	}
	return nil
}

func (m *SessionRepoMock) GetByToken(ctx context.Context, token string) (*entity.Session, error) {
	if m.GetByTokenFunc != nil {
		return m.GetByTokenFunc(ctx, token)
	}
	return nil, nil
}

func (m *SessionRepoMock) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Session, error) {
	if m.GetByUserIDFunc != nil {
		return m.GetByUserIDFunc(ctx, userID)
	}
	return nil, nil
}

func (m *SessionRepoMock) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *SessionRepoMock) DeleteByToken(ctx context.Context, token string) error {
	if m.DeleteByTokenFunc != nil {
		return m.DeleteByTokenFunc(ctx, token)
	}
	return nil
}
