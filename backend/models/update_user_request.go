package models

// UpdateUserRequest модель для обновления данных пользователя
type UpdateUserRequest struct {
	FirstName    *string  `json:"first_name"`
	LastName     *string  `json:"last_name"`
	Username     *string  `json:"username" validate:"omitempty,min=3,max=50"`
	Email        *string  `json:"email" validate:"omitempty,email"`
	MonthlyLimit *float64 `json:"monthly_limit" validate:"omitempty,gte=0"`
	SavingsGoal  *float64 `json:"savings_goal" validate:"omitempty,gte=0"`
}
