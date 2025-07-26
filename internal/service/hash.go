package service

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type hashService struct {
	logger *logrus.Logger
}

func NewHashService(logger *logrus.Logger) HashService {
	return &hashService{
		logger: logger,
	}
}

func (h *hashService) HashPassword(password string) (string, error) {
	h.logger.WithField("component", "hash_service").Debug("hashing password")

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.WithError(err).Error("failed to hash password")
		return "", err
	}

	h.logger.WithField("component", "hash_service").Debug("password hashed successfully")
	return string(bytes), nil
}

func (h *hashService) CheckPasswordHash(password, hash string) bool {
	h.logger.WithField("component", "hash_service").Debug("checking password hash")

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		h.logger.WithError(err).Debug("password hash check failed")
		return false
	}

	h.logger.WithField("component", "hash_service").Debug("password hash check successful")
	return true
}
