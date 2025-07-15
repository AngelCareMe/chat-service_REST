package entity

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Message struct {
	ID       uuid.UUID `json:"id"`
	Content  string    `json:"content"`
	SenderID uuid.UUID `json:"sender_id"`
	SentAt   time.Time `json:"sent_at"`
}

func (m *Message) Validate(content string) error {
	if content == "" {
		return fmt.Errorf("empty content")
	}
	return nil
}
