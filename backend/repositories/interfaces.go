package repositories

import (
	"context"
	"time"

	"cz.Finance/backend/models"
)

// UserRepository интерфейс для работы с пользователями в базе данных
type UserRepository interface {
	Create(ctx context.Context, user *models.User) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	UpdateAvatar(ctx context.Context, userID int64, avatarPath string) error
	Delete(ctx context.Context, id int64) error
}

// ExpenseRepository интерфейс для работы с тратами в базе данных
type ExpenseRepository interface {
	Create(ctx context.Context, expense *models.Expense) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.Expense, error)
	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]models.Expense, error)
	GetByUserIDAndPeriod(ctx context.Context, userID int64, startDate, endDate time.Time) ([]models.Expense, error)
	GetTotalAmountByUserIDAndPeriod(ctx context.Context, userID int64, startDate, endDate time.Time) (float64, error)
	GetCategorySummaryByUserIDAndPeriod(ctx context.Context, userID int64, startDate, endDate time.Time) (map[string]float64, error)
	Update(ctx context.Context, expense *models.Expense) error
	Delete(ctx context.Context, id int64, userID int64) error
}

// IncomeRepository интерфейс для работы с накоплениями в базе данных
type IncomeRepository interface {
	Create(ctx context.Context, income *models.Income) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.Income, error)
	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]models.Income, error)
	GetByUserIDAndPeriod(ctx context.Context, userID int64, startDate, endDate time.Time) ([]models.Income, error)
	GetTotalAmountByUserIDAndPeriod(ctx context.Context, userID int64, startDate, endDate time.Time) (float64, error)
	GetSourceSummaryByUserIDAndPeriod(ctx context.Context, userID int64, startDate, endDate time.Time) (map[string]float64, error)
	Update(ctx context.Context, income *models.Income) error
	Delete(ctx context.Context, id int64, userID int64) error
}

// WishlistRepository интерфейс для работы со списком желаний в базе данных
type WishlistRepository interface {
	Create(ctx context.Context, item *models.WishlistItem) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.WishlistItem, error)
	GetByUserID(ctx context.Context, userID int64) ([]models.WishlistItem, error)
	Update(ctx context.Context, item *models.WishlistItem) error
	Delete(ctx context.Context, id int64, userID int64) error
}

// TelegramUserRepository интерфейс для работы с Telegram пользователями
type TelegramUserRepository interface {
	Create(ctx context.Context, telegramUser *models.TelegramUser) (int64, error)
	GetByTelegramID(ctx context.Context, telegramID int64) (*models.TelegramUser, error)
	GetByUserID(ctx context.Context, userID int64) (*models.TelegramUser, error)
	Delete(ctx context.Context, id int64) error
}
