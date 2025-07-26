package usecase

import (
	"chat-service/internal/entity"
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type MessageRepository interface {
	Create(ctx context.Context, message *entity.Message) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Message, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Message, error)
	GetAll(ctx context.Context) ([]*entity.Message, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type SessionRepository interface {
	Create(ctx context.Context, session *entity.Session) error
	GetByToken(ctx context.Context, token string) (*entity.Session, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Session, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByToken(ctx context.Context, token string) error
}
