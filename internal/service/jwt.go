package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type jwtService struct {
	secretKey string
	logger    *logrus.Logger
}

func NewJWTService(secretKey string, logger *logrus.Logger) JWTService {
	return &jwtService{
		secretKey: secretKey,
		logger:    logger,
	}
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func (j *jwtService) GenerateToken(userID uuid.UUID) (string, error) {
	j.logger.WithField("user_id", userID).Debug("generating JWT token")

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.secretKey))

	if err != nil {
		j.logger.WithError(err).Error("failed to generate JWT token")
		return "", err
	}

	j.logger.WithField("user_id", userID).Debug("JWT token generated successfully")
	return signedToken, nil
}

func (j *jwtService) ValidateToken(tokenString string) (uuid.UUID, error) {
	j.logger.WithField("token", j.maskToken(tokenString)).Debug("validating JWT token")

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil {
		j.logger.WithError(err).Warn("failed to parse JWT token")
		return uuid.Nil, err
	}

	if !token.Valid {
		j.logger.Warn("invalid JWT token")
		return uuid.Nil, errors.New("invalid token")
	}

	j.logger.WithField("user_id", claims.UserID).Debug("JWT token validated successfully")
	return claims.UserID, nil
}

func (j *jwtService) maskToken(token string) string {
	if len(token) == 0 {
		return "<empty>"
	}
	if len(token) <= 8 {
		return "***"
	}
	// Берем минимум из длины токена и желаемой длины префикса (например, 10)
	end := 10
	if len(token) < end {
		end = len(token) // Это гарантирует, что end не превысит len(token)
	}
	return token[:end] + "..."
}
