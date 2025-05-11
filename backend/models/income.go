package models

import (
	"time"
)

// IncomeSource перечисляет возможные источники дохода
type IncomeSource string

const (
	SourceSalary     IncomeSource = "salary"
	SourceFreelance  IncomeSource = "freelance"
	SourceInvestment IncomeSource = "investment"
	SourceGift       IncomeSource = "gift"
	SourceRental     IncomeSource = "rental"
	SourceOther      IncomeSource = "other"
)

// Income представляет модель накопления (дохода)
type Income struct {
	ID          int64        `json:"id" db:"id"`
	UserID      int64        `json:"user_id" db:"user_id"`
	Amount      float64      `json:"amount" db:"amount" validate:"required,gt=0"`
	Source      IncomeSource `json:"source" db:"source" validate:"required"`
	Date        time.Time    `json:"date" db:"date"`
	Description string       `json:"description" db:"description"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
}

// CreateIncomeRequest модель для создания нового накопления
type CreateIncomeRequest struct {
	Amount      float64      `json:"amount" validate:"required,gt=0"`
	Source      IncomeSource `json:"source" validate:"required"`
	Date        time.Time    `json:"date"`
	Description string       `json:"description"`
}

// UpdateIncomeRequest модель для обновления накопления
type UpdateIncomeRequest struct {
	Amount      *float64      `json:"amount" validate:"omitempty,gt=0"`
	Source      *IncomeSource `json:"source"`
	Date        *time.Time    `json:"date"`
	Description *string       `json:"description"`
}

// IncomeSummary предоставляет общую информацию о накоплениях за период
type IncomeSummary struct {
	TotalAmount    float64            `json:"total_amount"`
	SavingsGoal    float64            `json:"savings_goal"`
	GoalPercentage float64            `json:"goal_percentage"`
	SourceSummary  map[string]float64 `json:"source_summary"`
	RecentIncomes  []Income           `json:"recent_incomes"`
}
