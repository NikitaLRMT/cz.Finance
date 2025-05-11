package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"cz.Finance/backend/models"
)

// PostgresUserRepository представляет реализацию репозитория пользователей на PostgreSQL
type PostgresUserRepository struct {
	db *sql.DB
}

// NewUserRepository создает новый экземпляр репозитория пользователей
func NewUserRepository(db *sql.DB) UserRepository {
	return &PostgresUserRepository{db: db}
}

// Create создает нового пользователя в базе данных
func (r *PostgresUserRepository) Create(ctx context.Context, user *models.User) (int64, error) {
	query := `
		INSERT INTO users (email, username, password_hash, first_name, last_name, monthly_limit, savings_goal, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.Username,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.MonthlyLimit,
		user.SavingsGoal,
		time.Now(),
		time.Now(),
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetByID получает пользователя по его ID
func (r *PostgresUserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	query := `
		SELECT id, email, username, password_hash, first_name, last_name, avatar_path, monthly_limit, savings_goal, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	// Создаем переменные для обработки возможных NULL значений
	var avatarPath sql.NullString
	var firstName sql.NullString
	var lastName sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&firstName,
		&lastName,
		&avatarPath,
		&user.MonthlyLimit,
		&user.SavingsGoal,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("пользователь не найден")
		}
		return nil, err
	}

	// Присваиваем значения из NullString
	if firstName.Valid {
		user.FirstName = firstName.String
	} else {
		user.FirstName = ""
	}

	if lastName.Valid {
		user.LastName = lastName.String
	} else {
		user.LastName = ""
	}

	if avatarPath.Valid {
		user.AvatarPath = avatarPath.String
	} else {
		user.AvatarPath = ""
	}

	return &user, nil
}

// GetByEmail получает пользователя по его электронной почте
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	fmt.Printf("REPO: Попытка получить пользователя по email: %s\n", email)

	query := `
		SELECT id, email, username, password_hash, first_name, last_name, avatar_path, monthly_limit, savings_goal, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user models.User
	// Создаем переменные для обработки возможных NULL значений
	var avatarPath sql.NullString
	var firstName sql.NullString
	var lastName sql.NullString

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&firstName,
		&lastName,
		&avatarPath,
		&user.MonthlyLimit,
		&user.SavingsGoal,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("REPO: Пользователь с email %s не найден\n", email)
			return nil, errors.New("пользователь не найден")
		}
		fmt.Printf("REPO: Ошибка при получении пользователя по email %s: %v\n", email, err)
		return nil, err
	}

	// Присваиваем значения из NullString
	if firstName.Valid {
		user.FirstName = firstName.String
	} else {
		user.FirstName = ""
	}

	if lastName.Valid {
		user.LastName = lastName.String
	} else {
		user.LastName = ""
	}

	if avatarPath.Valid {
		user.AvatarPath = avatarPath.String
	} else {
		user.AvatarPath = ""
	}

	fmt.Printf("REPO: Пользователь найден: ID=%d, Email=%s, Username=%s, PasswordHash=%s\n",
		user.ID, user.Email, user.Username, user.PasswordHash)
	return &user, nil
}

// GetByUsername получает пользователя по его имени пользователя
func (r *PostgresUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, email, username, password_hash, first_name, last_name, avatar_path, monthly_limit, savings_goal, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	var user models.User
	// Создаем переменные для обработки возможных NULL значений
	var avatarPath sql.NullString
	var firstName sql.NullString
	var lastName sql.NullString

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&firstName,
		&lastName,
		&avatarPath,
		&user.MonthlyLimit,
		&user.SavingsGoal,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("пользователь не найден")
		}
		return nil, err
	}

	// Присваиваем значения из NullString
	if firstName.Valid {
		user.FirstName = firstName.String
	} else {
		user.FirstName = ""
	}

	if lastName.Valid {
		user.LastName = lastName.String
	} else {
		user.LastName = ""
	}

	if avatarPath.Valid {
		user.AvatarPath = avatarPath.String
	} else {
		user.AvatarPath = ""
	}

	return &user, nil
}

// Update обновляет данные пользователя в базе данных
func (r *PostgresUserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET email = $1, username = $2, first_name = $3, last_name = $4, monthly_limit = $5, savings_goal = $6, updated_at = $7
		WHERE id = $8
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.Email,
		user.Username,
		user.FirstName,
		user.LastName,
		user.MonthlyLimit,
		user.SavingsGoal,
		time.Now(),
		user.ID,
	)

	return err
}

// UpdateAvatar обновляет путь к аватару пользователя
func (r *PostgresUserRepository) UpdateAvatar(ctx context.Context, userID int64, avatarPath string) error {
	// Добавляем логирование перед обновлением
	fmt.Printf("Обновление аватара пользователя ID=%d: %s\n", userID, avatarPath)

	query := `
		UPDATE users
		SET avatar_path = $1, updated_at = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		avatarPath,
		time.Now(),
		userID,
	)

	if err != nil {
		fmt.Printf("Ошибка при обновлении аватара: %v\n", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("пользователь не найден")
	}

	fmt.Printf("Аватар успешно обновлен для пользователя ID=%d\n", userID)
	return nil
}

// Delete удаляет пользователя из базы данных
func (r *PostgresUserRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("пользователь не найден")
	}

	return nil
}
