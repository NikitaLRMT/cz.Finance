package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"cz.Finance/backend/models"
)

// PostgresTelegramUserRepository представляет реализацию репозитория Telegram пользователей на PostgreSQL
type PostgresTelegramUserRepository struct {
	db *sql.DB
}

// NewTelegramUserRepository создает новый экземпляр репозитория Telegram пользователей
func NewTelegramUserRepository(db *sql.DB) TelegramUserRepository {
	return &PostgresTelegramUserRepository{db: db}
}

// Create создает новую связь Telegram пользователя с пользователем системы
func (r *PostgresTelegramUserRepository) Create(ctx context.Context, telegramUser *models.TelegramUser) (int64, error) {
	query := `
		INSERT INTO telegram_users (user_id, telegram_id, username, first_name, last_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		telegramUser.UserID,
		telegramUser.TelegramID,
		telegramUser.Username,
		telegramUser.FirstName,
		telegramUser.LastName,
		time.Now(),
		time.Now(),
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetByTelegramID получает связь по Telegram ID
func (r *PostgresTelegramUserRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*models.TelegramUser, error) {
	query := `
		SELECT id, user_id, telegram_id, username, first_name, last_name
		FROM telegram_users
		WHERE telegram_id = $1
	`

	var telegramUser models.TelegramUser
	err := r.db.QueryRowContext(ctx, query, telegramID).Scan(
		&telegramUser.ID,
		&telegramUser.UserID,
		&telegramUser.TelegramID,
		&telegramUser.Username,
		&telegramUser.FirstName,
		&telegramUser.LastName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("связь с Telegram пользователем не найдена")
		}
		return nil, err
	}

	return &telegramUser, nil
}

// GetByUserID получает связь по ID пользователя системы
func (r *PostgresTelegramUserRepository) GetByUserID(ctx context.Context, userID int64) (*models.TelegramUser, error) {
	query := `
		SELECT id, user_id, telegram_id, username, first_name, last_name
		FROM telegram_users
		WHERE user_id = $1
	`

	var telegramUser models.TelegramUser
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&telegramUser.ID,
		&telegramUser.UserID,
		&telegramUser.TelegramID,
		&telegramUser.Username,
		&telegramUser.FirstName,
		&telegramUser.LastName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("связь с Telegram пользователем не найдена")
		}
		return nil, err
	}

	return &telegramUser, nil
}

// Delete удаляет связь с Telegram пользователем
func (r *PostgresTelegramUserRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM telegram_users WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
