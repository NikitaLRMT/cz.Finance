package models

import (
	"time"
)

// ExpenseCategory перечисляет возможные категории трат
type ExpenseCategory string

const (
	CategoryFood          ExpenseCategory = "food"
	CategoryTransport     ExpenseCategory = "transport"
	CategoryHousing       ExpenseCategory = "housing"
	CategoryUtilities     ExpenseCategory = "utilities"
	CategoryShopping      ExpenseCategory = "shopping"
	CategoryEntertainment ExpenseCategory = "entertainment"
	CategoryHealthcare    ExpenseCategory = "healthcare"
	CategoryEducation     ExpenseCategory = "education"
	CategoryTravel        ExpenseCategory = "travel"
	CategoryOther         ExpenseCategory = "other"
)

// Expense представляет модель траты
type Expense struct {
	ID          int64           `json:"id" db:"id"`
	UserID      int64           `json:"user_id" db:"user_id"`
	Title       string          `json:"title" db:"title" validate:"required,min=2,max=100"`
	Amount      float64         `json:"amount" db:"amount" validate:"required,gt=0"`
	Category    ExpenseCategory `json:"category" db:"category" validate:"required"`
	Date        time.Time       `json:"date" db:"date"`
	Description string          `json:"description" db:"description"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

// CreateExpenseRequest модель для создания новой траты
type CreateExpenseRequest struct {
	Title       string          `json:"title" validate:"required,min=2,max=100"`
	Amount      float64         `json:"amount" validate:"required,gt=0"`
	Category    ExpenseCategory `json:"category" validate:"required"`
	Date        time.Time       `json:"date"`
	Description string          `json:"description"`
}

// UpdateExpenseRequest модель для обновления траты
type UpdateExpenseRequest struct {
	Title       *string          `json:"title" validate:"omitempty,min=2,max=100"`
	Amount      *float64         `json:"amount" validate:"omitempty,gt=0"`
	Category    *ExpenseCategory `json:"category"`
	Date        *time.Time       `json:"date"`
	Description *string          `json:"description"`
}

// ExpenseSummary предоставляет общую информацию о тратах за период
type ExpenseSummary struct {
	TotalAmount     float64            `json:"total_amount"`
	MonthlyLimit    float64            `json:"monthly_limit"`
	LimitPercentage float64            `json:"limit_percentage"`
	CategorySummary map[string]float64 `json:"category_summary"`
	RecentExpenses  []Expense          `json:"recent_expenses"`
}
