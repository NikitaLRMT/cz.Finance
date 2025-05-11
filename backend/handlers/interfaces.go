package handlers

import (
	"net/http"
)

// UserHandler интерфейс для обработки запросов связанных с пользователями
type UserHandler interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	UploadAvatar(w http.ResponseWriter, r *http.Request)
	RemoveAvatar(w http.ResponseWriter, r *http.Request)
}

// ExpenseHandler интерфейс для обработки запросов связанных с тратами
type ExpenseHandler interface {
	CreateExpense(w http.ResponseWriter, r *http.Request)
	GetExpense(w http.ResponseWriter, r *http.Request)
	GetUserExpenses(w http.ResponseWriter, r *http.Request)
	GetExpenseSummary(w http.ResponseWriter, r *http.Request)
	UpdateExpense(w http.ResponseWriter, r *http.Request)
	DeleteExpense(w http.ResponseWriter, r *http.Request)
}

// IncomeHandler интерфейс для обработки запросов связанных с накоплениями
type IncomeHandler interface {
	CreateIncome(w http.ResponseWriter, r *http.Request)
	GetIncome(w http.ResponseWriter, r *http.Request)
	GetUserIncomes(w http.ResponseWriter, r *http.Request)
	GetIncomeSummary(w http.ResponseWriter, r *http.Request)
	UpdateIncome(w http.ResponseWriter, r *http.Request)
	DeleteIncome(w http.ResponseWriter, r *http.Request)
}

// DashboardHandler интерфейс для обработки запросов связанных с панелью статистики
type DashboardHandler interface {
	GetDashboardSummary(w http.ResponseWriter, r *http.Request)
	GetMonthlyStats(w http.ResponseWriter, r *http.Request)
	GetYearlyStats(w http.ResponseWriter, r *http.Request)
	GetBudgetGoals(w http.ResponseWriter, r *http.Request)
	SetBudgetGoal(w http.ResponseWriter, r *http.Request)
}

// CalculatorHandler интерфейс для обработки запросов связанных с калькуляторами
type CalculatorHandler interface {
	CompoundInterestCalculator(w http.ResponseWriter, r *http.Request)
	MortgageCalculator(w http.ResponseWriter, r *http.Request)
}

// WishlistHandler интерфейс для обработки запросов связанных со списком желаний
type WishlistHandler interface {
	CreateWishlistItem(w http.ResponseWriter, r *http.Request)
	GetWishlistItem(w http.ResponseWriter, r *http.Request)
	GetUserWishlist(w http.ResponseWriter, r *http.Request)
	UpdateWishlistItem(w http.ResponseWriter, r *http.Request)
	DeleteWishlistItem(w http.ResponseWriter, r *http.Request)
}

// TelegramHandler интерфейс для обработки запросов от Telegram
type TelegramHandler interface {
	LinkTelegramAccount(w http.ResponseWriter, r *http.Request)
	GetUserByTelegramID(w http.ResponseWriter, r *http.Request)
	UnlinkTelegramAccount(w http.ResponseWriter, r *http.Request)
}
