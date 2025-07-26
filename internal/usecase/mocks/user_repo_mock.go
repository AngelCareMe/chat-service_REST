package mocks

import (
	"context"

	"chat-service/internal/entity"

	"github.com/google/uuid"
)

type UserRepoMock struct {
	CreateFunc     func(ctx context.Context, user *entity.User) error
	GetByIDFunc    func(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmailFunc func(ctx context.Context, email string) (*entity.User, error)
	UpdateFunc     func(ctx context.Context, user *entity.User) error
	DeleteFunc     func(ctx context.Context, id uuid.UUID) error
}

func (m *UserRepoMock) Create(ctx context.Context, user *entity.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return nil
}

func (m *UserRepoMock) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *UserRepoMock) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (m *UserRepoMock) Update(ctx context.Context, user *entity.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, user)
	}
	return nil
}

func (m *UserRepoMock) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}
