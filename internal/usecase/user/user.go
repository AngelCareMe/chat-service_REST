package usecase

import (
	"chat-service/internal/entity"
	"chat-service/internal/usecase/auth"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

type UserUseCase struct {
	userRepo    UserRepository
	authService auth.AuthService
	logger      *logrus.Logger
}

func NewUserUseCase(userRepo UserRepository, authService auth.AuthService, logger *logrus.Logger) *UserUseCase {
	return &UserUseCase{
		userRepo:    userRepo,
		authService: authService,
		logger:      logger,
	}
}

func (uc *UserUseCase) Register(ctx context.Context, username, password string) (*entity.User, string, error) {
	if username == "" || password == "" {
		return nil, "", fmt.Errorf("username or password is empty")
	}

	if _, err := uc.userRepo.GetByUsername(ctx, username); err == nil {
		return nil, "", fmt.Errorf("username already exists")
	}

	user := &entity.User{
		Username:  username,
		Password:  password,
		CreatedAt: time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, "", fmt.Errorf("create user: %w", err)
	}

	if err := user.Validate(); err != nil {
		return nil, "", fmt.Errorf("validate user: %w", err)
	}

	token, err := uc.authService.GenerateJWT(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("generate jwt: %w", err)
	}

	uc.logger.WithFields(logrus.Fields{
		"username": username,
		"user_id":  user.ID,
	}).Info("User registered")

	return user, token, nil
}

func (uc *UserUseCase) Login(ctx context.Context, username, password string) (*entity.User, string, error) {
	user, err := uc.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, "", fmt.Errorf("get user: %w", err)
	}

	if err := uc.authService.VerifyPassword(password); err != nil {
		return nil, "", fmt.Errorf("verify password: %w", err)
	}

	token, err := uc.authService.GenerateJWT(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("generate jwt: %w", err)
	}

	uc.logger.WithFields(logrus.Fields{
		"username": username,
		"user_id":  user.ID,
	}).Info("User logged in")

	return user, token, nil
}

func (uc *UserUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	uc.logger.WithFields(logrus.Fields{
		"user_id": user.ID,
	}).Info("User fetched")

	return user, nil
}

func (uc *UserUseCase) Update(ctx context.Context, id uuid.UUID, username, password string) error {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	if username != "" && username != user.Username {
		if _, err := uc.userRepo.GetByUsername(ctx, username); err != nil {
			return fmt.Errorf("username %s already exists", username)
		}
		user.Username = username
	}

	if password != "" {
		user.Password = password
	}

	if err := user.Validate(); err != nil {
		return fmt.Errorf("validate user: %w", err)
	}

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	uc.logger.WithFields(logrus.Fields{
		"user_id": user.ID,
	}).Info("User updated")

	return nil
}

func (uc *UserUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := uc.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	uc.logger.WithFields(logrus.Fields{
		"user_id": id,
	}).Info("User deleted")

	return nil
}
