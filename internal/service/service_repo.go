package service

import "github.com/google/uuid"

type HashService interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type JWTService interface {
	GenerateToken(userID uuid.UUID) (string, error)
	ValidateToken(token string) (uuid.UUID, error)
}
