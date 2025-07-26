package handler

import (
	"net/http"
	"strings"

	"chat-service/internal/entity"
	"chat-service/internal/usecase/session"
	"chat-service/internal/usecase/user"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	userUsecase    user.UserUsecase
	sessionUsecase session.SessionUsecase
	logger         *logrus.Logger
}

func NewUserHandler(
	userUsecase user.UserUsecase,
	sessionUsecase session.SessionUsecase,
	logger *logrus.Logger,
) *UserHandler {
	return &UserHandler{
		userUsecase:    userUsecase,
		sessionUsecase: sessionUsecase,
		logger:         logger,
	}
}

// RegisterRequest структура для регистрации
// swagger:model RegisterRequest
type RegisterRequest struct {
	// Username пользователя
	// required: true
	// min length: 3
	Username string `json:"username" binding:"required,min=3"`

	// Email пользователя
	// required: true
	// format: email
	Email string `json:"email" binding:"required,email"`

	// Пароль пользователя
	// required: true
	// min length: 6
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest структура для логина
// swagger:model LoginRequest
type LoginRequest struct {
	// Email пользователя
	// required: true
	// format: email
	Email string `json:"email" binding:"required,email"`

	// Пароль пользователя
	// required: true
	Password string `json:"password" binding:"required"`
}

// UserResponse структура ответа с пользователем
// swagger:model UserResponse
type UserResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    *entity.User `json:"data"`
}

// Register регистрирует нового пользователя
// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя в системе
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body RegisterRequest true "Данные для регистрации"
// @Success 201 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Warn("invalid register request body")
		SendError(c, "Invalid request", err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.WithField("email", req.Email).Info("user registration attempt")

	user, err := h.userUsecase.Register(c.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		h.logger.WithError(err).Warn("user registration failed")
		HandleError(c, err, h.logger)
		return
	}

	// Создаем сессию для нового пользователя
	session, err := h.sessionUsecase.CreateSession(c.Request.Context(), user.ID)
	if err != nil {
		h.logger.WithError(err).Error("failed to create session after registration")
		SendError(c, "Registration successful but login failed", "Please login manually", http.StatusOK)
		return
	}

	response := struct {
		User    *entity.User    `json:"user"`
		Session *entity.Session `json:"session"`
	}{
		User:    user,
		Session: session,
	}

	h.logger.WithField("user_id", user.ID).Info("user registered successfully")
	SendSuccess(c, response, "User registered successfully", http.StatusCreated)
}

// Login аутентифицирует пользователя
// @Summary Вход в систему
// @Description Аутентифицирует пользователя и возвращает токен
// @Tags users
// @Accept  json
// @Produce  json
// @Param credentials body LoginRequest true "Учетные данные"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Warn("invalid login request body")
		SendError(c, "Invalid request", err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.WithField("email", req.Email).Info("user login attempt")

	user, err := h.userUsecase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.WithError(err).Warn("user login failed")
		HandleError(c, err, h.logger)
		return
	}

	// Создаем сессию
	session, err := h.sessionUsecase.CreateSession(c.Request.Context(), user.ID)
	if err != nil {
		h.logger.WithError(err).Error("failed to create session after login")
		SendError(c, "Login failed", "Failed to create session", http.StatusInternalServerError)
		return
	}

	response := struct {
		User    *entity.User    `json:"user"`
		Session *entity.Session `json:"session"`
	}{
		User:    user,
		Session: session,
	}

	h.logger.WithField("user_id", user.ID).Info("user logged in successfully")
	SendSuccess(c, response, "Login successful", http.StatusOK)
}

// GetProfile возвращает профиль текущего пользователя
// @Summary Получение профиля пользователя
// @Description Возвращает профиль авторизованного пользователя
// @Tags users
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {object} UserResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, err := GetUserFromContext(c)
	if err != nil {
		h.logger.WithError(err).Warn("failed to get user from context")
		HandleError(c, err, h.logger)
		return
	}

	h.logger.WithField("user_id", userID).Debug("fetching user profile")

	user, err := h.userUsecase.GetProfile(c.Request.Context(), userID)
	if err != nil {
		h.logger.WithError(err).Error("failed to fetch user profile")
		HandleError(c, err, h.logger)
		return
	}

	h.logger.WithField("user_id", userID).Debug("user profile fetched successfully")
	SendSuccess(c, user, "Profile retrieved successfully", http.StatusOK)
}

// Logout выходит из аккаунта
// @Summary Выход из системы
// @Description Завершает сессию пользователя
// @Tags users
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Router /logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	userID, err := GetUserFromContext(c)
	if err != nil {
		h.logger.WithError(err).Warn("failed to get user from context")
		HandleError(c, err, h.logger)
		return
	}

	// Получаем токен из заголовка
	authHeader := c.GetHeader("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	if tokenString != authHeader && tokenString != "" {
		// Удаляем сессию
		err = h.sessionUsecase.DeleteSession(c.Request.Context(), tokenString)
		if err != nil {
			h.logger.WithError(err).Warn("failed to delete session")
		}
	}

	h.logger.WithField("user_id", userID).Info("user logged out successfully")
	SendSuccess(c, nil, "Logged out successfully", http.StatusOK)
}

// UpdateProfile обновляет профиль пользователя
// @Summary Обновление профиля пользователя
// @Description Обновляет данные профиля авторизованного пользователя
// @Tags users
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param user body object{username=string,email=string} false "Данные для обновления"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, err := GetUserFromContext(c)
	if err != nil {
		h.logger.WithError(err).Warn("failed to get user from context")
		HandleError(c, err, h.logger)
		return
	}

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Warn("invalid update profile request body")
		SendError(c, "Invalid request", err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем текущего пользователя
	user, err := h.userUsecase.GetProfile(c.Request.Context(), userID)
	if err != nil {
		h.logger.WithError(err).Error("failed to fetch user for update")
		HandleError(c, err, h.logger)
		return
	}

	// Обновляем только переданные поля
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	err = h.userUsecase.UpdateProfile(c.Request.Context(), user)
	if err != nil {
		h.logger.WithError(err).Error("failed to update user profile")
		HandleError(c, err, h.logger)
		return
	}

	h.logger.WithField("user_id", userID).Info("user profile updated successfully")
	SendSuccess(c, user, "Profile updated successfully", http.StatusOK)
}

// DeleteUser удаляет аккаунт пользователя
// @Summary Удаление аккаунта пользователя
// @Description Удаляет аккаунт авторизованного пользователя
// @Tags users
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /profile [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := GetUserFromContext(c)
	if err != nil {
		h.logger.WithError(err).Warn("failed to get user from context")
		HandleError(c, err, h.logger)
		return
	}

	h.logger.WithField("user_id", userID).Warn("user account deletion requested")

	err = h.userUsecase.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		h.logger.WithError(err).Error("failed to delete user account")
		HandleError(c, err, h.logger)
		return
	}

	h.logger.WithField("user_id", userID).Info("user account deleted successfully")
	SendSuccess(c, nil, "Account deleted successfully", http.StatusOK)
}
