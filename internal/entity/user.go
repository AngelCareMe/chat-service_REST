package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) Validate() error {
	if u.Username == "" {
		return &ValidationError{"username is required"}
	}
	if len(u.Username) < 3 {
		return &ValidationError{"username must be at least 3 characters"}
	}
	if u.Email == "" {
		return &ValidationError{"email is required"}
	}
	if u.Password == "" {
		return &ValidationError{"password is required"}
	}
	if len(u.Password) < 6 {
		return &ValidationError{"password must be at least 6 characters"}
	}
	return nil
}

// ValidateForUpdate метод для валидации при обновлении
func (u *User) ValidateForUpdate() error {
	if u.Username == "" {
		return &ValidationError{"username is required"}
	}
	if len(u.Username) < 3 {
		return &ValidationError{"username must be at least 3 characters"}
	}
	if u.Email == "" {
		return &ValidationError{"email is required"}
	}
	// Пароль не обязателен при обновлении
	return nil
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
