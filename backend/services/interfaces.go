package services

import (
	"context"
	"mime/multipart"
	"time"

	"cz.Finance/backend/models"
)

// UserService интерфейс для работы с пользователями
type UserService interface {
	SignUp(ctx context.Context, signup *models.UserSignup) (*models.TokenResponse, error)
	Login(ctx context.Context, login *models.UserLogin) (*models.TokenResponse, error)
	GetUser(ctx context.Context, id int64) (*models.UserResponse, error)
	UpdateUser(ctx context.Context, id int64, updateRequest *models.UpdateUserRequest) (*models.UserResponse, error)
	UploadAvatar(ctx context.Context, userID int64, file *multipart.FileHeader) (*models.UserResponse, error)
	RemoveAvatar(ctx context.Context, userID int64) (*models.UserResponse, error)
	DeleteUser(ctx context.Context, id int64) error
}

// ExpenseService интерфейс для работы с тратами
type ExpenseService interface {
	CreateExpense(ctx context.Context, userID int64, request *models.CreateExpenseRequest) (*models.Expense, error)
	GetExpense(ctx context.Context, id int64, userID int64) (*models.Expense, error)
	GetUserExpenses(ctx context.Context, userID int64, limit, offset int) ([]models.Expense, error)
	GetUserExpensesByPeriod(ctx context.Context, userID int64, startDate, endDate time.Time) ([]models.Expense, error)
	GetExpenseSummary(ctx context.Context, userID int64, startDate, endDate time.Time) (*models.ExpenseSummary, error)
	UpdateExpense(ctx context.Context, id int64, userID int64, request *models.UpdateExpenseRequest) (*models.Expense, error)
	DeleteExpense(ctx context.Context, id int64, userID int64) error
}

// IncomeService интерфейс для работы с накоплениями
type IncomeService interface {
	CreateIncome(ctx context.Context, userID int64, request *models.CreateIncomeRequest) (*models.Income, error)
	GetIncome(ctx context.Context, id int64, userID int64) (*models.Income, error)
	GetUserIncomes(ctx context.Context, userID int64, limit, offset int) ([]models.Income, error)
	GetUserIncomesByPeriod(ctx context.Context, userID int64, startDate, endDate time.Time) ([]models.Income, error)
	GetIncomeSummary(ctx context.Context, userID int64, startDate, endDate time.Time) (*models.IncomeSummary, error)
	UpdateIncome(ctx context.Context, id int64, userID int64, request *models.UpdateIncomeRequest) (*models.Income, error)
	DeleteIncome(ctx context.Context, id int64, userID int64) error
}

// AuthService интерфейс для аутентификации
type AuthService interface {
	GenerateToken(userID int64, email string) (string, time.Time, error)
	ValidateToken(token string) (*models.TokenClaims, error)
}

// DashboardService интерфейс для статистики и информационной панели
type DashboardService interface {
	GetDashboardSummary(ctx context.Context, userID int64, limit int) (map[string]interface{}, error)
	GetMonthlyStats(ctx context.Context, userID int64, year int, month int) (map[string]interface{}, error)
	GetYearlyStats(ctx context.Context, userID int64, year int) (map[string]interface{}, error)
	GetBudgetGoals(ctx context.Context, userID int64) (map[string]interface{}, error)
	SetBudgetGoal(ctx context.Context, userID int64, category string, amount float64) error
}

// WishlistService интерфейс для работы со списком желаний
type WishlistService interface {
	CreateWishlistItem(ctx context.Context, userID int64, request *models.CreateWishlistItemRequest) (*models.WishlistItem, error)
	GetWishlistItem(ctx context.Context, id int64, userID int64) (*models.WishlistItem, error)
	GetUserWishlist(ctx context.Context, userID int64) ([]models.WishlistItem, error)
	UpdateWishlistItem(ctx context.Context, id int64, userID int64, request *models.UpdateWishlistItemRequest) (*models.WishlistItem, error)
	DeleteWishlistItem(ctx context.Context, id int64, userID int64) error
}

// TelegramService интерфейс для работы с Telegram пользователями
type TelegramService interface {
	LinkAccount(ctx context.Context, telegramID int64, username, firstName, lastName string, request *models.TelegramLinkRequest) (*models.User, error)
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*models.User, error)
	UnlinkAccount(ctx context.Context, telegramID int64) error
}
