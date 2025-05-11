package models

import (
	"fmt"
	"time"
)

// User представляет модель пользователя в системе
type User struct {
	ID           int64     `json:"id" db:"id"`
	Email        string    `json:"email" db:"email" validate:"required,email"`
	Username     string    `json:"username" db:"username" validate:"required,min=3,max=50"`
	PasswordHash string    `json:"-" db:"password_hash"`
	FirstName    string    `json:"first_name" db:"first_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	AvatarPath   string    `json:"avatar_url" db:"avatar_path"`
	MonthlyLimit float64   `json:"monthly_limit" db:"monthly_limit"`
	SavingsGoal  float64   `json:"savings_goal" db:"savings_goal"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// UserSignup модель для регистрации пользователя
type UserSignup struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// UserLogin модель для входа пользователя
type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse модель для ответа пользователю без чувствительных данных
type UserResponse struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	AvatarURL    string    `json:"avatar_url"`
	MonthlyLimit float64   `json:"monthly_limit"`
	SavingsGoal  float64   `json:"savings_goal"`
	CreatedAt    time.Time `json:"created_at"`
}

// ToUserResponse конвертирует User в UserResponse
func (u *User) ToUserResponse() UserResponse {
	// Явное логирование AvatarPath перед конвертацией
	fmt.Printf("█████████ КОНВЕРТАЦИЯ █████████\n")
	fmt.Printf("User.ID: %d\n", u.ID)
	fmt.Printf("User.AvatarPath: '%s'\n", u.AvatarPath)

	resp := UserResponse{
		ID:           u.ID,
		Email:        u.Email,
		Username:     u.Username,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		AvatarURL:    u.AvatarPath,
		MonthlyLimit: u.MonthlyLimit,
		SavingsGoal:  u.SavingsGoal,
		CreatedAt:    u.CreatedAt,
	}

	// Проверка на корректность конвертации
	fmt.Printf("UserResponse.AvatarURL: '%s'\n", resp.AvatarURL)
	fmt.Printf("█████████ КОНЕЦ КОНВЕРТАЦИИ █████████\n")

	return resp
}
