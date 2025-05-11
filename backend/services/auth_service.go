package services

import (
	"time"

	"cz.Finance/backend/configs"
	"cz.Finance/backend/models"
	"cz.Finance/backend/utils"
)

// AuthServiceImpl представляет реализацию сервиса аутентификации
type AuthServiceImpl struct {
	jwtConfig configs.JWTConfig
}

// NewAuthService создает новый экземпляр сервиса аутентификации
func NewAuthService(jwtConfig configs.JWTConfig) AuthService {
	return &AuthServiceImpl{
		jwtConfig: jwtConfig,
	}
}

// GenerateToken создает JWT токен для пользователя
func (s *AuthServiceImpl) GenerateToken(userID int64, email string) (string, time.Time, error) {
	return utils.GenerateJWT(userID, email, s.jwtConfig.Secret, s.jwtConfig.ExpiresIn)
}

// ValidateToken проверяет JWT токен и возвращает информацию о пользователе
func (s *AuthServiceImpl) ValidateToken(token string) (*models.TokenClaims, error) {
	return utils.ValidateJWT(token, s.jwtConfig.Secret)
}
