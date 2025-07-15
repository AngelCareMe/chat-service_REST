package usecase

import (
	"chat-service/internal/entity"
	"context"
	"github.com/google/uuid"
)

type UserUsecase interface {
	Register(ctx context.Context, username, password string) (*entity.User, string, error)
	Login(ctx context.Context, username, password string) (*entity.User, string, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	Update(ctx context.Context, id uuid.UUID, username, password string) error
	Delete(ctx context.Context, id uuid.UUID) error
}
