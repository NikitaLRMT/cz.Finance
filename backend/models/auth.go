package models

import (
	"time"
)

// TokenClaims представляет структуру для JWT токена
type TokenClaims struct {
	UserID int64  `json:"id"`
	Email  string `json:"email,omitempty"`
}

// TokenResponse модель с токеном аутентификации
type TokenResponse struct {
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expires_at"`
	User      UserResponse `json:"user"`
}

// ErrorResponse стандартная модель для ответа с ошибкой
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}
