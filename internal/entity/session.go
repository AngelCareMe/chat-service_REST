package entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Session) Validate() error {
	if s.UserID == uuid.Nil {
		return &ValidationError{"user_id is required"}
	}
	if s.Token == "" {
		return &ValidationError{"token is required"}
	}
	if s.ExpiresAt.IsZero() {
		return &ValidationError{"expires_at is required"}
	}
	if s.ExpiresAt.Before(time.Now()) {
		return &ValidationError{"token is expired"}
	}
	return nil
}
