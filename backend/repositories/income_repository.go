package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"cz.Finance/backend/models"
)

// PostgresIncomeRepository представляет реализацию репозитория накоплений на PostgreSQL
type PostgresIncomeRepository struct {
	db *sql.DB
}

// NewIncomeRepository создает новый экземпляр репозитория накоплений
func NewIncomeRepository(db *sql.DB) IncomeRepository {
	return &PostgresIncomeRepository{db: db}
}

// Create создает новое накопление в базе данных
func (r *PostgresIncomeRepository) Create(ctx context.Context, income *models.Income) (int64, error) {
	query := `
		INSERT INTO incomes (user_id, amount, source, date, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		income.UserID,
		income.Amount,
		income.Source,
		income.Date,
		income.Description,
		time.Now(),
		time.Now(),
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetByID получает накопление по его ID
func (r *PostgresIncomeRepository) GetByID(ctx context.Context, id int64) (*models.Income, error) {
	query := `
		SELECT id, user_id, amount, source, date, description, created_at, updated_at
		FROM incomes
		WHERE id = $1
	`

	var income models.Income
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&income.ID,
		&income.UserID,
		&income.Amount,
		&income.Source,
		&income.Date,
		&income.Description,
		&income.CreatedAt,
		&income.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("накопление не найдено")
		}
		return nil, err
	}

	return &income, nil
}

// GetByUserID получает накопления пользователя с пагинацией
func (r *PostgresIncomeRepository) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]models.Income, error) {
	query := `
		SELECT id, user_id, amount, source, date, description, created_at, updated_at
		FROM incomes
		WHERE user_id = $1
		ORDER BY date DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incomes []models.Income
	for rows.Next() {
		var income models.Income
		err := rows.Scan(
			&income.ID,
			&income.UserID,
			&income.Amount,
			&income.Source,
			&income.Date,
			&income.Description,
			&income.CreatedAt,
			&income.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		incomes = append(incomes, income)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return incomes, nil
}

// GetByUserIDAndPeriod получает накопления пользователя за определенный период
func (r *PostgresIncomeRepository) GetByUserIDAndPeriod(ctx context.Context, userID int64, startDate, endDate time.Time) ([]models.Income, error) {
	query := `
		SELECT id, user_id, amount, source, date, description, created_at, updated_at
		FROM incomes
		WHERE user_id = $1 AND date >= $2 AND date <= $3
		ORDER BY date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incomes []models.Income
	for rows.Next() {
		var income models.Income
		err := rows.Scan(
			&income.ID,
			&income.UserID,
			&income.Amount,
			&income.Source,
			&income.Date,
			&income.Description,
			&income.CreatedAt,
			&income.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		incomes = append(incomes, income)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return incomes, nil
}

// GetTotalAmountByUserIDAndPeriod получает общую сумму накоплений пользователя за период
func (r *PostgresIncomeRepository) GetTotalAmountByUserIDAndPeriod(ctx context.Context, userID int64, startDate, endDate time.Time) (float64, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM incomes
		WHERE user_id = $1 AND date >= $2 AND date <= $3
	`

	var totalAmount float64
	err := r.db.QueryRowContext(ctx, query, userID, startDate, endDate).Scan(&totalAmount)
	if err != nil {
		return 0, err
	}

	return totalAmount, nil
}

// GetSourceSummaryByUserIDAndPeriod получает сводку накоплений по источникам за период
func (r *PostgresIncomeRepository) GetSourceSummaryByUserIDAndPeriod(ctx context.Context, userID int64, startDate, endDate time.Time) (map[string]float64, error) {
	query := `
		SELECT source, SUM(amount) as total
		FROM incomes
		WHERE user_id = $1 AND date >= $2 AND date <= $3
		GROUP BY source
		ORDER BY total DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sourceSummary := make(map[string]float64)
	for rows.Next() {
		var source string
		var amount float64
		err := rows.Scan(&source, &amount)
		if err != nil {
			return nil, err
		}
		sourceSummary[source] = amount
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sourceSummary, nil
}

// Update обновляет информацию о накоплении
func (r *PostgresIncomeRepository) Update(ctx context.Context, income *models.Income) error {
	query := `
		UPDATE incomes
		SET amount = $1, source = $2, date = $3, description = $4, updated_at = $5
		WHERE id = $6 AND user_id = $7
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		income.Amount,
		income.Source,
		income.Date,
		income.Description,
		time.Now(),
		income.ID,
		income.UserID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("накопление не найдено или у вас нет прав на его изменение")
	}

	return nil
}

// Delete удаляет накопление из базы данных
func (r *PostgresIncomeRepository) Delete(ctx context.Context, id int64, userID int64) error {
	query := `DELETE FROM incomes WHERE id = $1 AND user_id = $2`

	result, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("накопление не найдено или у вас нет прав на его удаление")
	}

	return nil
}
