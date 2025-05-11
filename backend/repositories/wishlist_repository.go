package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"cz.Finance/backend/models"
)

// PostgresWishlistRepository представляет реализацию репозитория списка желаний на PostgreSQL
type PostgresWishlistRepository struct {
	db *sql.DB
}

// NewWishlistRepository создает новый экземпляр репозитория списка желаний
func NewWishlistRepository(db *sql.DB) WishlistRepository {
	return &PostgresWishlistRepository{db: db}
}

// Create создает новый элемент списка желаний в базе данных
func (r *PostgresWishlistRepository) Create(ctx context.Context, item *models.WishlistItem) (int64, error) {
	query := `
		INSERT INTO wishlist (user_id, title, price, priority, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		item.UserID,
		item.Title,
		item.Price,
		item.Priority,
		item.Description,
		time.Now(),
		time.Now(),
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetByID получает элемент списка желаний по его ID
func (r *PostgresWishlistRepository) GetByID(ctx context.Context, id int64) (*models.WishlistItem, error) {
	query := `
		SELECT id, user_id, title, price, priority, description, created_at, updated_at
		FROM wishlist
		WHERE id = $1
	`

	var item models.WishlistItem
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.UserID,
		&item.Title,
		&item.Price,
		&item.Priority,
		&item.Description,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("элемент списка желаний не найден")
		}
		return nil, err
	}

	return &item, nil
}

// GetByUserID получает все элементы списка желаний пользователя
func (r *PostgresWishlistRepository) GetByUserID(ctx context.Context, userID int64) ([]models.WishlistItem, error) {
	query := `
		SELECT id, user_id, title, price, priority, description, created_at, updated_at
		FROM wishlist
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.WishlistItem
	for rows.Next() {
		var item models.WishlistItem
		err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Title,
			&item.Price,
			&item.Priority,
			&item.Description,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// Update обновляет элемент списка желаний в базе данных
func (r *PostgresWishlistRepository) Update(ctx context.Context, item *models.WishlistItem) error {
	query := `
		UPDATE wishlist
		SET title = $1, price = $2, priority = $3, description = $4, updated_at = $5
		WHERE id = $6 AND user_id = $7
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		item.Title,
		item.Price,
		item.Priority,
		item.Description,
		time.Now(),
		item.ID,
		item.UserID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("элемент списка желаний не найден или не принадлежит пользователю")
	}

	return nil
}

// Delete удаляет элемент списка желаний из базы данных
func (r *PostgresWishlistRepository) Delete(ctx context.Context, id int64, userID int64) error {
	query := `
		DELETE FROM wishlist
		WHERE id = $1 AND user_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("элемент списка желаний не найден или не принадлежит пользователю")
	}

	return nil
}
