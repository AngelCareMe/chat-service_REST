package session

import (
	"chat-service/internal/entity"
	"chat-service/internal/service"
	"chat-service/internal/usecase"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type sessionUsecase struct {
	sessionRepo usecase.SessionRepository
	jwtService  service.JWTService
	logger      *logrus.Logger
}

func NewSessionUsecase(sessionRepo usecase.SessionRepository, jwtService service.JWTService, logger *logrus.Logger) SessionUsecase {
	return &sessionUsecase{
		sessionRepo: sessionRepo,
		jwtService:  jwtService,
		logger:      logger,
	}
}

func (s *sessionUsecase) CreateSession(ctx context.Context, userID uuid.UUID) (*entity.Session, error) {
	s.logger.WithField("user_id", userID).Info("creating new session")

	// Генерируем JWT токен
	s.logger.WithField("user_id", userID).Debug("generating JWT token")
	token, err := s.jwtService.GenerateToken(userID)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("failed to generate JWT token")
		return nil, err
	}

	session := &entity.Session{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 часа
		CreatedAt: time.Now(),
	}

	if err := session.Validate(); err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Warn("session validation failed")
		return nil, err
	}

	s.logger.WithField("session_id", session.ID).Debug("saving session to repository")
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		s.logger.WithError(err).WithField("session_id", session.ID).Error("failed to create session")
		return nil, err
	}

	s.logger.WithField("session_id", session.ID).Info("session created successfully")
	return session, nil
}

func (s *sessionUsecase) ValidateSession(ctx context.Context, token string) (*entity.Session, error) {
	s.logger.WithField("token", token[:min(20, len(token))]+"...").Debug("validating session")

	session, err := s.sessionRepo.GetByToken(ctx, token)
	if err != nil {
		s.logger.WithField("token", token[:min(20, len(token))]+"...").Warn("session not found")
		return nil, &BusinessError{"invalid session"}
	}

	if session.ExpiresAt.Before(time.Now()) {
		// Удаляем просроченную сессию
		s.logger.WithField("session_id", session.ID).Warn("session expired, cleaning up")
		s.sessionRepo.DeleteByToken(ctx, token)
		return nil, &BusinessError{"session expired"}
	}

	// Проверяем JWT токен
	s.logger.Debug("validating JWT token")
	userID, err := s.jwtService.ValidateToken(token)
	if err != nil {
		s.logger.WithError(err).WithField("token", token[:min(20, len(token))]+"...").Warn("invalid JWT token")
		return nil, &BusinessError{"invalid token"}
	}

	if userID != session.UserID {
		s.logger.WithFields(logrus.Fields{
			"expected_user_id": session.UserID,
			"actual_user_id":   userID,
		}).Warn("token user ID mismatch")
		return nil, &BusinessError{"token mismatch"}
	}

	s.logger.WithField("session_id", session.ID).Debug("session validated successfully")
	return session, nil
}

func (s *sessionUsecase) DeleteSession(ctx context.Context, token string) error {
	s.logger.WithField("token", token[:min(20, len(token))]+"...").Warn("deleting session")

	err := s.sessionRepo.DeleteByToken(ctx, token)
	if err != nil {
		s.logger.WithError(err).WithField("token", token[:min(20, len(token))]+"...").Error("failed to delete session")
		return err
	}

	s.logger.WithField("token", token[:min(20, len(token))]+"...").Info("session deleted successfully")
	return nil
}

type BusinessError struct {
	Message string
}

func (e *BusinessError) Error() string {
	return e.Message
}

func (e *BusinessError) ValidationError() bool {
	return true
}
