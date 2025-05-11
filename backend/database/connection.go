package database

import (
	"database/sql"
	"fmt"
	"os"

	"cz.Finance/backend/configs"

	_ "github.com/lib/pq" // Драйвер PostgreSQL
)

// ConnectDB устанавливает соединение с базой данных PostgreSQL
func ConnectDB(config *configs.Config) (*sql.DB, error) {
	var dsn string

	// Проверяем наличие переменной окружения DATABASE_URL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL != "" {
		dsn = dbURL
	} else {
		// Если DATABASE_URL не задан, строим DSN из отдельных настроек
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			config.Database.Host,
			config.Database.Port,
			config.Database.User,
			config.Database.Password,
			config.Database.DBName,
			config.Database.SSLMode)
	}

	// Открываем соединение с базой данных
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	// Проверяем соединение
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка проверки соединения с базой данных: %w", err)
	}

	return db, nil
}
