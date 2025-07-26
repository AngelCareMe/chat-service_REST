package user

import (
	"chat-service/internal/entity"
	"context"

	"github.com/google/uuid"
)

type UserUsecase interface {
	Register(ctx context.Context, username, email, password string) (*entity.User, error)
	Login(ctx context.Context, email, password string) (*entity.User, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	UpdateProfile(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}
