package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"cz.Finance/backend/models"
)

// PostgresExpenseRepository представляет реализацию репозитория трат на PostgreSQL
type PostgresExpenseRepository struct {
	db *sql.DB
}

// NewExpenseRepository создает новый экземпляр репозитория трат
func NewExpenseRepository(db *sql.DB) ExpenseRepository {
	return &PostgresExpenseRepository{db: db}
}

// Create создает новую трату в базе данных
func (r *PostgresExpenseRepository) Create(ctx context.Context, expense *models.Expense) (int64, error) {
	query := `
		INSERT INTO expenses (user_id, title, amount, category, date, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		expense.UserID,
		expense.Title,
		expense.Amount,
		expense.Category,
		expense.Date,
		expense.Description,
		time.Now(),
		time.Now(),
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetByID получает трату по её ID
func (r *PostgresExpenseRepository) GetByID(ctx context.Context, id int64) (*models.Expense, error) {
	query := `
		SELECT id, user_id, title, amount, category, date, description, created_at, updated_at
		FROM expenses
		WHERE id = $1
	`

	var expense models.Expense
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&expense.ID,
		&expense.UserID,
		&expense.Title,
		&expense.Amount,
		&expense.Category,
		&expense.Date,
		&expense.Description,
		&expense.CreatedAt,
		&expense.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("трата не найдена")
		}
		return nil, err
	}

	return &expense, nil
}

// GetByUserID получает траты пользователя с пагинацией
func (r *PostgresExpenseRepository) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]models.Expense, error) {
	query := `
		SELECT id, user_id, title, amount, category, date, description, created_at, updated_at
		FROM expenses
		WHERE user_id = $1
		ORDER BY date DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var expense models.Expense
		err := rows.Scan(
			&expense.ID,
			&expense.UserID,
			&expense.Title,
			&expense.Amount,
			&expense.Category,
			&expense.Date,
			&expense.Description,
			&expense.CreatedAt,
			&expense.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}

// GetByUserIDAndPeriod получает траты пользователя за определенный период
func (r *PostgresExpenseRepository) GetByUserIDAndPeriod(ctx context.Context, userID int64, startDate, endDate time.Time) ([]models.Expense, error) {
	query := `
		SELECT id, user_id, title, amount, category, date, description, created_at, updated_at
		FROM expenses
		WHERE user_id = $1 AND date >= $2 AND date <= $3
		ORDER BY date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var expense models.Expense
		err := rows.Scan(
			&expense.ID,
			&expense.UserID,
			&expense.Title,
			&expense.Amount,
			&expense.Category,
			&expense.Date,
			&expense.Description,
			&expense.CreatedAt,
			&expense.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}

// GetTotalAmountByUserIDAndPeriod получает общую сумму трат пользователя за период
func (r *PostgresExpenseRepository) GetTotalAmountByUserIDAndPeriod(ctx context.Context, userID int64, startDate, endDate time.Time) (float64, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM expenses
		WHERE user_id = $1 AND date >= $2 AND date <= $3
	`

	var totalAmount float64
	err := r.db.QueryRowContext(ctx, query, userID, startDate, endDate).Scan(&totalAmount)
	if err != nil {
		return 0, err
	}

	return totalAmount, nil
}

// GetCategorySummaryByUserIDAndPeriod получает сводку трат по категориям за период
func (r *PostgresExpenseRepository) GetCategorySummaryByUserIDAndPeriod(ctx context.Context, userID int64, startDate, endDate time.Time) (map[string]float64, error) {
	query := `
		SELECT category, SUM(amount) as total
		FROM expenses
		WHERE user_id = $1 AND date >= $2 AND date <= $3
		GROUP BY category
		ORDER BY total DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categorySummary := make(map[string]float64)
	for rows.Next() {
		var category string
		var amount float64
		err := rows.Scan(&category, &amount)
		if err != nil {
			return nil, err
		}
		categorySummary[category] = amount
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categorySummary, nil
}

// Update обновляет информацию о трате
func (r *PostgresExpenseRepository) Update(ctx context.Context, expense *models.Expense) error {
	query := `
		UPDATE expenses
		SET title = $1, amount = $2, category = $3, date = $4, description = $5, updated_at = $6
		WHERE id = $7 AND user_id = $8
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		expense.Title,
		expense.Amount,
		expense.Category,
		expense.Date,
		expense.Description,
		time.Now(),
		expense.ID,
		expense.UserID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("трата не найдена или у вас нет прав на её изменение")
	}

	return nil
}

// Delete удаляет трату из базы данных
func (r *PostgresExpenseRepository) Delete(ctx context.Context, id int64, userID int64) error {
	query := `DELETE FROM expenses WHERE id = $1 AND user_id = $2`

	result, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("трата не найдена или у вас нет прав на её удаление")
	}

	return nil
}
