package message

import (
	"chat-service/internal/entity"
	"context"

	"github.com/google/uuid"
)

type MessageUsecase interface {
	CreateMessage(ctx context.Context, userID uuid.UUID, content string) (*entity.Message, error)
	GetMessageByID(ctx context.Context, messageID uuid.UUID) (*entity.Message, error)
	GetMessagesByUser(ctx context.Context, userID uuid.UUID) ([]*entity.Message, error)
	GetAllMessages(ctx context.Context) ([]*entity.Message, error)
	DeleteMessage(ctx context.Context, messageID uuid.UUID) error
}
