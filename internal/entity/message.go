package entity

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Message) Validate() error {
	if m.UserID == uuid.Nil {
		return &ValidationError{"user_id is required"}
	}
	if m.Content == "" {
		return &ValidationError{"content is required"}
	}
	if len(m.Content) > 1000 {
		return &ValidationError{"content must be less than 1000 characters"}
	}
	return nil
}
