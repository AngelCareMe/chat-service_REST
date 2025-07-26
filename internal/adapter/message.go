package postgres

import (
	"context"
	"fmt"

	"chat-service/internal/entity"
	"chat-service/internal/usecase"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type messageRepo struct {
	adapter *PostgresAdapter
	psql    squirrel.StatementBuilderType
}

func NewMessageRepository(adapter *PostgresAdapter) usecase.MessageRepository {
	return &messageRepo{
		adapter: adapter,
		psql:    squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *messageRepo) Create(ctx context.Context, message *entity.Message) error {
	// Валидация перед вставкой
	if err := r.validateMessage(message); err != nil {
		return err
	}

	query, args, err := r.psql.Insert("messages").
		Columns("id", "user_id", "content", "created_at", "updated_at").
		Values(message.ID, message.UserID, message.Content, message.CreatedAt, message.UpdatedAt).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		r.adapter.logger.WithError(err).Error("failed to build insert query for message")
		return fmt.Errorf("failed to build query: %w", err)
	}

	var returnedID uuid.UUID
	err = r.adapter.QueryRow(ctx, query, args...).Scan(&returnedID)
	if err != nil {
		r.adapter.logger.WithError(err).WithField("message_id", message.ID).Error("failed to create message in database")
		return fmt.Errorf("failed to insert message: %w", err)
	}

	r.adapter.logger.WithField("message_id", returnedID).Info("message created successfully in database")
	return nil
}

func (r *messageRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Message, error) {
	if id == uuid.Nil {
		return nil, &ValidationError{"invalid message ID"}
	}

	query, args, err := r.psql.Select("id", "user_id", "content", "created_at", "updated_at").
		From("messages").
		Where(squirrel.Eq{"id": id}).
		Limit(1).
		ToSql()

	if err != nil {
		r.adapter.logger.WithError(err).Error("failed to build select query for message by ID")
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var message entity.Message
	err = r.adapter.QueryRow(ctx, query, args...).Scan(
		&message.ID, &message.UserID, &message.Content, &message.CreatedAt, &message.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			r.adapter.logger.WithField("message_id", id).Warn("message not found")
			return nil, &NotFoundError{"message not found"}
		}
		r.adapter.logger.WithError(err).WithField("message_id", id).Error("failed to get message by ID")
		return nil, fmt.Errorf("failed to query message: %w", err)
	}

	r.adapter.logger.WithField("message_id", message.ID).Debug("message retrieved by ID")
	return &message, nil
}

func (r *messageRepo) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Message, error) {
	if userID == uuid.Nil {
		return nil, &ValidationError{"invalid user ID"}
	}

	query, args, err := r.psql.Select("id", "user_id", "content", "created_at", "updated_at").
		From("messages").
		Where(squirrel.Eq{"user_id": userID}).
		OrderBy("created_at DESC").
		ToSql()

	if err != nil {
		r.adapter.logger.WithError(err).Error("failed to build select query for messages by user ID")
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.adapter.Query(ctx, query, args...)
	if err != nil {
		r.adapter.logger.WithError(err).WithField("user_id", userID).Error("failed to query messages by user ID")
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []*entity.Message
	for rows.Next() {
		var message entity.Message
		err := rows.Scan(&message.ID, &message.UserID, &message.Content, &message.CreatedAt, &message.UpdatedAt)
		if err != nil {
			r.adapter.logger.WithError(err).WithField("user_id", userID).Error("failed to scan message row")
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, &message)
	}

	// Проверяем ошибки при итерации
	if err = rows.Err(); err != nil {
		r.adapter.logger.WithError(err).WithField("user_id", userID).Error("error during message rows iteration")
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	r.adapter.logger.WithField("user_id", userID).Debugf("retrieved %d messages for user", len(messages))
	return messages, nil
}

func (r *messageRepo) GetAll(ctx context.Context) ([]*entity.Message, error) {
	query, args, err := r.psql.Select("id", "user_id", "content", "created_at", "updated_at").
		From("messages").
		OrderBy("created_at DESC").
		ToSql()

	if err != nil {
		r.adapter.logger.WithError(err).Error("failed to build select query for all messages")
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.adapter.Query(ctx, query, args...)
	if err != nil {
		r.adapter.logger.WithError(err).Error("failed to query all messages")
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []*entity.Message
	for rows.Next() {
		var message entity.Message
		err := rows.Scan(&message.ID, &message.UserID, &message.Content, &message.CreatedAt, &message.UpdatedAt)
		if err != nil {
			r.adapter.logger.WithError(err).Error("failed to scan message row")
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, &message)
	}

	// Проверяем ошибки при итерации
	if err = rows.Err(); err != nil {
		r.adapter.logger.WithError(err).Error("error during all messages rows iteration")
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	r.adapter.logger.Debugf("retrieved %d messages total", len(messages))
	return messages, nil
}

func (r *messageRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return &ValidationError{"invalid message ID"}
	}

	query, args, err := r.psql.Delete("messages").
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		r.adapter.logger.WithError(err).Error("failed to build delete query for message")
		return fmt.Errorf("failed to build query: %w", err)
	}

	var deletedID uuid.UUID
	err = r.adapter.QueryRow(ctx, query, args...).Scan(&deletedID)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.adapter.logger.WithField("message_id", id).Warn("message not found for deletion")
			return &NotFoundError{"message not found"}
		}
		r.adapter.logger.WithError(err).WithField("message_id", id).Error("failed to delete message")
		return fmt.Errorf("failed to delete message: %w", err)
	}

	r.adapter.logger.WithField("message_id", deletedID).Info("message deleted successfully")
	return nil
}

// Валидация сообщения
func (r *messageRepo) validateMessage(message *entity.Message) error {
	if message == nil {
		return &ValidationError{"message cannot be nil"}
	}

	if message.UserID == uuid.Nil {
		return &ValidationError{"user_id is required"}
	}

	if message.Content == "" {
		return &ValidationError{"content is required"}
	}

	if len(message.Content) > 1000 {
		return &ValidationError{"content must be less than 1000 characters"}
	}

	return nil
}
