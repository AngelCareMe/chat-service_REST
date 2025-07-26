package user

import (
	"chat-service/internal/entity"
	"chat-service/internal/service"
	"chat-service/internal/usecase"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type userUsecase struct {
	userRepo    usecase.UserRepository
	sessionRepo usecase.SessionRepository
	hashService service.HashService
	jwtService  service.JWTService
	logger      *logrus.Logger
}

func NewUserUsecase(
	userRepo usecase.UserRepository,
	sessionRepo usecase.SessionRepository,
	hashService service.HashService,
	jwtService service.JWTService,
	logger *logrus.Logger,
) UserUsecase {
	return &userUsecase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		hashService: hashService,
		jwtService:  jwtService,
		logger:      logger,
	}
}

func (u *userUsecase) Register(ctx context.Context, username, email, password string) (*entity.User, error) {
	u.logger.WithFields(logrus.Fields{
		"username": username,
		"email":    email,
	}).Info("registering new user")

	// Проверяем, существует ли пользователь с таким email
	existingUser, _ := u.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		err := &BusinessError{"user with this email already exists"}
		u.logger.WithField("email", email).Warn("user already exists")
		return nil, err
	}

	// Хэшируем пароль
	u.logger.Debug("hashing user password")
	hashedPassword, err := u.hashService.HashPassword(password)
	if err != nil {
		u.logger.WithError(err).Error("failed to hash password")
		return nil, err
	}

	user := &entity.User{
		ID:        uuid.New(),
		Username:  username,
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Валидируем пользователя
	if err := user.Validate(); err != nil {
		u.logger.WithError(err).Warn("user validation failed")
		return nil, err
	}

	// Создаем пользователя
	u.logger.WithField("user_id", user.ID).Debug("creating user in repository")
	if err := u.userRepo.Create(ctx, user); err != nil {
		u.logger.WithError(err).WithField("user_id", user.ID).Error("failed to create user")
		return nil, err
	}

	// Очищаем пароль перед возвратом
	user.Password = ""
	u.logger.WithField("user_id", user.ID).Info("user registered successfully")
	return user, nil
}

func (u *userUsecase) Login(ctx context.Context, email, password string) (*entity.User, error) {
	u.logger.WithField("email", email).Info("user login attempt")

	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		u.logger.WithField("email", email).Warn("user not found during login")
		return nil, &BusinessError{"invalid credentials"}
	}

	// Проверяем пароль
	u.logger.Debug("checking password hash")
	if !u.hashService.CheckPasswordHash(password, user.Password) {
		u.logger.WithField("email", email).Warn("invalid password during login")
		return nil, &BusinessError{"invalid credentials"}
	}

	// Очищаем пароль перед возвратом
	user.Password = ""
	u.logger.WithField("user_id", user.ID).Info("user login successful")
	return user, nil
}

func (u *userUsecase) GetProfile(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	u.logger.WithField("user_id", userID).Debug("fetching user profile")

	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		u.logger.WithError(err).WithField("user_id", userID).Error("failed to fetch user profile")
		return nil, err
	}

	user.Password = "" // Очищаем пароль перед возвратом
	u.logger.WithField("user_id", userID).Debug("user profile fetched successfully")
	return user, nil
}

func (u *userUsecase) UpdateProfile(ctx context.Context, user *entity.User) error {
	u.logger.WithField("user_id", user.ID).Info("updating user profile")

	user.UpdatedAt = time.Now()
	if err := user.ValidateForUpdate(); err != nil {
		u.logger.WithError(err).WithField("user_id", user.ID).Warn("user validation failed during update")
		return err
	}

	err := u.userRepo.Update(ctx, user)
	if err != nil {
		u.logger.WithError(err).WithField("user_id", user.ID).Error("failed to update user profile")
		return err
	}

	u.logger.WithField("user_id", user.ID).Info("user profile updated successfully")
	return nil
}

func (u *userUsecase) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	u.logger.WithField("user_id", userID).Warn("deleting user")

	// Удаляем сессии пользователя
	_, err := u.sessionRepo.GetByUserID(ctx, userID)
	if err == nil {
		// Здесь можно добавить удаление всех сессий пользователя
		u.logger.WithField("user_id", userID).Debug("cleaning up user sessions")
	}

	err = u.userRepo.Delete(ctx, userID)
	if err != nil {
		u.logger.WithError(err).WithField("user_id", userID).Error("failed to delete user")
		return err
	}

	u.logger.WithField("user_id", userID).Info("user deleted successfully")
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
