package postgres

import (
	"context"
	"fmt"
	"time"

	"chat-service/internal/entity"
	"chat-service/internal/usecase"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type sessionRepo struct {
	adapter *PostgresAdapter
	psql    squirrel.StatementBuilderType
}

func NewSessionRepository(adapter *PostgresAdapter) usecase.SessionRepository {
	return &sessionRepo{
		adapter: adapter,
		psql:    squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *sessionRepo) Create(ctx context.Context, session *entity.Session) error {
	// Валидация перед вставкой
	if err := r.validateSession(session); err != nil {
		return err
	}

	query, args, err := r.psql.Insert("sessions").
		Columns("id", "user_id", "token", "expires_at", "created_at").
		Values(session.ID, session.UserID, session.Token, session.ExpiresAt, session.CreatedAt).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		r.adapter.logger.WithError(err).Error("failed to build insert query for session")
		return fmt.Errorf("failed to build query: %w", err)
	}

	var returnedID uuid.UUID
	err = r.adapter.QueryRow(ctx, query, args...).Scan(&returnedID)
	if err != nil {
		r.adapter.logger.WithError(err).WithField("session_id", session.ID).Error("failed to create session in database")
		return fmt.Errorf("failed to insert session: %w", err)
	}

	r.adapter.logger.WithField("session_id", returnedID).Info("session created successfully in database")
	return nil
}

func (r *sessionRepo) GetByToken(ctx context.Context, token string) (*entity.Session, error) {
	if token == "" {
		return nil, &ValidationError{"token is required"}
	}

	query, args, err := r.psql.Select("id", "user_id", "token", "expires_at", "created_at").
		From("sessions").
		Where(squirrel.Eq{"token": token}).
		Limit(1).
		ToSql()

	if err != nil {
		r.adapter.logger.WithError(err).Error("failed to build select query for session by token")
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var session entity.Session
	err = r.adapter.QueryRow(ctx, query, args...).Scan(
		&session.ID, &session.UserID, &session.Token, &session.ExpiresAt, &session.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			r.adapter.logger.WithField("token", r.maskToken(token)).Warn("session not found by token")
			return nil, &NotFoundError{"session not found"}
		}
		r.adapter.logger.WithError(err).WithField("token", r.maskToken(token)).Error("failed to get session by token")
		return nil, fmt.Errorf("failed to query session: %w", err)
	}

	// Проверяем срок действия
	if session.ExpiresAt.Before(time.Now()) {
		r.adapter.logger.WithField("session_id", session.ID).Warn("session expired")
		return nil, &ValidationError{"session expired"}
	}

	r.adapter.logger.WithField("session_id", session.ID).Debug("session retrieved by token")
	return &session, nil
}

func (r *sessionRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Session, error) {
	if userID == uuid.Nil {
		return nil, &ValidationError{"invalid user ID"}
	}

	query, args, err := r.psql.Select("id", "user_id", "token", "expires_at", "created_at").
		From("sessions").
		Where(squirrel.Eq{"user_id": userID}).
		OrderBy("created_at DESC").
		Limit(1).
		ToSql()

	if err != nil {
		r.adapter.logger.WithError(err).Error("failed to build select query for session by user ID")
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var session entity.Session
	err = r.adapter.QueryRow(ctx, query, args...).Scan(
		&session.ID, &session.UserID, &session.Token, &session.ExpiresAt, &session.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			r.adapter.logger.WithField("user_id", userID).Warn("session not found by user ID")
			return nil, &NotFoundError{"session not found"}
		}
		r.adapter.logger.WithError(err).WithField("user_id", userID).Error("failed to get session by user ID")
		return nil, fmt.Errorf("failed to query session: %w", err)
	}

	// Проверяем срок действия
	if session.ExpiresAt.Before(time.Now()) {
		r.adapter.logger.WithField("session_id", session.ID).Warn("session expired")
		return nil, &ValidationError{"session expired"}
	}

	r.adapter.logger.WithField("session_id", session.ID).Debug("session retrieved by user ID")
	return &session, nil
}

func (r *sessionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return &ValidationError{"invalid session ID"}
	}

	query, args, err := r.psql.Delete("sessions").
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		r.adapter.logger.WithError(err).Error("failed to build delete query for session")
		return fmt.Errorf("failed to build query: %w", err)
	}

	var deletedID uuid.UUID
	err = r.adapter.QueryRow(ctx, query, args...).Scan(&deletedID)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.adapter.logger.WithField("session_id", id).Warn("session not found for deletion")
			return &NotFoundError{"session not found"}
		}
		r.adapter.logger.WithError(err).WithField("session_id", id).Error("failed to delete session")
		return fmt.Errorf("failed to delete session: %w", err)
	}

	r.adapter.logger.WithField("session_id", deletedID).Info("session deleted successfully")
	return nil
}

func (r *sessionRepo) DeleteByToken(ctx context.Context, token string) error {
	if token == "" {
		return &ValidationError{"token is required"}
	}

	query, args, err := r.psql.Delete("sessions").
		Where(squirrel.Eq{"token": token}).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		r.adapter.logger.WithError(err).Error("failed to build delete query for session by token")
		return fmt.Errorf("failed to build query: %w", err)
	}

	var deletedID uuid.UUID
	err = r.adapter.QueryRow(ctx, query, args...).Scan(&deletedID)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.adapter.logger.WithField("token", r.maskToken(token)).Warn("session not found for deletion by token")
			return &NotFoundError{"session not found"}
		}
		r.adapter.logger.WithError(err).WithField("token", r.maskToken(token)).Error("failed to delete session by token")
		return fmt.Errorf("failed to delete session: %w", err)
	}

	r.adapter.logger.WithField("session_id", deletedID).Info("session deleted successfully by token")
	return nil
}

// Валидация сессии
func (r *sessionRepo) validateSession(session *entity.Session) error {
	if session == nil {
		return &ValidationError{"session cannot be nil"}
	}

	if session.UserID == uuid.Nil {
		return &ValidationError{"user_id is required"}
	}

	if session.Token == "" {
		return &ValidationError{"token is required"}
	}

	if session.ExpiresAt.IsZero() {
		return &ValidationError{"expires_at is required"}
	}

	if session.ExpiresAt.Before(time.Now()) {
		return &ValidationError{"cannot create expired session"}
	}

	return nil
}

// Маскировка токена для логов
func (r *sessionRepo) maskToken(token string) string {
	if len(token) <= 8 {
		return "***"
	}
	return token[:4] + "..." + token[len(token)-4:]
}
