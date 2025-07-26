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

type userRepo struct {
	adapter *PostgresAdapter
	psql    squirrel.StatementBuilderType
}

func NewUserRepository(adapter *PostgresAdapter) usecase.UserRepository {
	return &userRepo{
		adapter: adapter,
		psql:    squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *userRepo) Create(ctx context.Context, user *entity.User) error {
	// Валидация перед вставкой
	if err := r.validateUser(user); err != nil {
		return err
	}

	query, args, err := r.psql.Insert("users").
		Columns("id", "username", "email", "password", "created_at", "updated_at").
		Values(user.ID, user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		r.adapter.logger.WithError(err).Error("failed to build insert query for user")
		return fmt.Errorf("failed to build query: %w", err)
	}

	var returnedID uuid.UUID
	err = r.adapter.QueryRow(ctx, query, args...).Scan(&returnedID)
	if err != nil {
		r.adapter.logger.WithError(err).Error("failed to create user in database")
		return fmt.Errorf("failed to insert user: %w", err)
	}

	r.adapter.logger.WithField("user_id", returnedID).Info("user created successfully in database")
	return nil
}

func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	if id == uuid.Nil {
		return nil, &ValidationError{"invalid user ID"}
	}

	query, args, err := r.psql.Select("id", "username", "email", "password", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"id": id}).
		Limit(1).
		ToSql()

	if err != nil {
		r.adapter.logger.WithError(err).Error("failed to build select query for user by ID")
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var user entity.User
	err = r.adapter.QueryRow(ctx, query, args...).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			r.adapter.logger.WithField("user_id", id).Warn("user not found")
			return nil, &NotFoundError{"user not found"}
		}
		r.adapter.logger.WithError(err).WithField("user_id", id).Error("failed to get user by ID")
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	r.adapter.logger.WithField("user_id", user.ID).Debug("user retrieved by ID")
	return &user, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	if email == "" {
		return nil, &ValidationError{"email is required"}
	}

	query, args, err := r.psql.Select("id", "username", "email", "password", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"email": email}).
		Limit(1).
		ToSql()

	if err != nil {
		r.adapter.logger.WithError(err).Error("failed to build select query for user by email")
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var user entity.User
	err = r.adapter.QueryRow(ctx, query, args...).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			r.adapter.logger.WithField("email", email).Warn("user not found by email")
			return nil, &NotFoundError{"user not found"}
		}
		r.adapter.logger.WithError(err).WithField("email", email).Error("failed to get user by email")
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	r.adapter.logger.WithField("user_id", user.ID).Debug("user retrieved by email")
	return &user, nil
}

func (r *userRepo) Update(ctx context.Context, user *entity.User) error {
	// При обновлении не проверяем пароль, если он пустой
	queryBuilder := r.psql.Update("users").
		Set("username", user.Username).
		Set("email", user.Email).
		Set("updated_at", user.UpdatedAt).
		Where(squirrel.Eq{"id": user.ID}).
		Suffix("RETURNING id")

	// Если пароль не пустой, обновляем его
	if user.Password != "" {
		queryBuilder = queryBuilder.Set("password", user.Password)
	}

	query, args, err := queryBuilder.ToSql()

	if err != nil {
		r.adapter.logger.WithError(err).Error("failed to build update query for user")
		return fmt.Errorf("failed to build query: %w", err)
	}

	var returnedID uuid.UUID
	err = r.adapter.QueryRow(ctx, query, args...).Scan(&returnedID)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.adapter.logger.WithField("user_id", user.ID).Warn("user not found for update")
			return &NotFoundError{"user not found"}
		}
		r.adapter.logger.WithError(err).WithField("user_id", user.ID).Error("failed to update user")
		return fmt.Errorf("failed to update user: %w", err)
	}

	r.adapter.logger.WithField("user_id", returnedID).Info("user updated successfully")
	return nil
}

func (r *userRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return &ValidationError{"invalid user ID"}
	}

	// Атомарная операция: удаляем пользователя и все связанные данные
	tx, err := r.adapter.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			r.adapter.logger.WithField("user_id", id).Warn("transaction rolled back")
		}
	}()

	// Удаляем все сообщения пользователя
	msgQuery, msgArgs, err := r.psql.Delete("messages").
		Where(squirrel.Eq{"user_id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete messages query: %w", err)
	}

	err = r.adapter.ExecTx(ctx, tx, msgQuery, msgArgs...)
	if err != nil {
		r.adapter.logger.WithError(err).WithField("user_id", id).Error("failed to delete user messages")
		return fmt.Errorf("failed to delete user messages: %w", err)
	}

	// Удаляем все сессии пользователя
	sessQuery, sessArgs, err := r.psql.Delete("sessions").
		Where(squirrel.Eq{"user_id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete sessions query: %w", err)
	}

	err = r.adapter.ExecTx(ctx, tx, sessQuery, sessArgs...)
	if err != nil {
		r.adapter.logger.WithError(err).WithField("user_id", id).Error("failed to delete user sessions")
		return fmt.Errorf("failed to delete user sessions: %w", err)
	}

	// Удаляем пользователя
	userQuery, userArgs, err := r.psql.Delete("users").
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete user query: %w", err)
	}

	var deletedID uuid.UUID
	err = r.adapter.QueryRowTx(ctx, tx, userQuery, userArgs...).Scan(&deletedID)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.adapter.logger.WithField("user_id", id).Warn("user not found for deletion")
			return &NotFoundError{"user not found"}
		}
		r.adapter.logger.WithError(err).WithField("user_id", id).Error("failed to delete user")
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Коммитим транзакцию
	err = tx.Commit(ctx)
	if err != nil {
		r.adapter.logger.WithError(err).WithField("user_id", id).Error("failed to commit transaction")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.adapter.logger.WithField("user_id", deletedID).Info("user and all related data deleted successfully")
	return nil
}

// Валидация пользователя
func (r *userRepo) validateUser(user *entity.User) error {
	if user == nil {
		return &ValidationError{"user cannot be nil"}
	}

	if user.Username == "" {
		return &ValidationError{"username is required"}
	}

	if len(user.Username) < 3 {
		return &ValidationError{"username must be at least 3 characters"}
	}

	if user.Email == "" {
		return &ValidationError{"email is required"}
	}

	if user.Password == "" {
		return &ValidationError{"password is required"}
	}

	if len(user.Password) < 6 {
		return &ValidationError{"password must be at least 6 characters"}
	}

	return nil
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}
