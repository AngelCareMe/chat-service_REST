package session

import (
	"chat-service/internal/entity"
	"context"

	"github.com/google/uuid"
)

type SessionUsecase interface {
	CreateSession(ctx context.Context, userID uuid.UUID) (*entity.Session, error)
	ValidateSession(ctx context.Context, token string) (*entity.Session, error)
	DeleteSession(ctx context.Context, token string) error
}
