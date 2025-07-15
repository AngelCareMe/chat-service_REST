package entity

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) Validate() error {
	if u.Username == "" || u.Password == "" {
		return fmt.Errorf("password or username can't be empty")
	}
	if len(u.Username) < 4 {
		return fmt.Errorf("username is too short")
	}
	if len(u.Password) < 5 {
		return fmt.Errorf("password is too short")
	}
	return nil
}
