package auth

import "github.com/google/uuid"

type AuthService interface {
	GeneratePasswordHash(password string) (string, error)
	VerifyPassword(password string) error
	GenerateJWT(userID uuid.UUID) (string, error)
	ValidateJWT(tokenString string) (int64, error)
}
