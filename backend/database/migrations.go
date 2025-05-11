package database

import (
	"database/sql"
	"fmt"
	"log"
)

// Migrations содержит все SQL-запросы для создания схемы базы данных
var Migrations = []string{
	`
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    avatar_path TEXT,
    monthly_limit DECIMAL(12, 2) DEFAULT 0,
    savings_goal DECIMAL(12, 2) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
`,
	`
CREATE TABLE IF NOT EXISTS expenses (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(100) NOT NULL,
    amount DECIMAL(12, 2) NOT NULL,
    category VARCHAR(50) NOT NULL,
    date TIMESTAMP WITH TIME ZONE DEFAULT now(),
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
`,
	`
CREATE TABLE IF NOT EXISTS incomes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    amount DECIMAL(12, 2) NOT NULL,
    source VARCHAR(50) NOT NULL,
    date TIMESTAMP WITH TIME ZONE DEFAULT now(),
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
`,
	`
CREATE TABLE IF NOT EXISTS wishlist (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(100) NOT NULL,
    price DECIMAL(12, 2) NOT NULL,
    priority VARCHAR(20) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
`,
	`
CREATE INDEX IF NOT EXISTS idx_expenses_user_id ON expenses(user_id);
CREATE INDEX IF NOT EXISTS idx_expenses_date ON expenses(date);
CREATE INDEX IF NOT EXISTS idx_expenses_category ON expenses(category);
CREATE INDEX IF NOT EXISTS idx_incomes_user_id ON incomes(user_id);
CREATE INDEX IF NOT EXISTS idx_incomes_date ON incomes(date);
CREATE INDEX IF NOT EXISTS idx_incomes_source ON incomes(source);
CREATE INDEX IF NOT EXISTS idx_wishlist_user_id ON wishlist(user_id);
CREATE INDEX IF NOT EXISTS idx_wishlist_priority ON wishlist(priority);
`,
	`
CREATE TABLE IF NOT EXISTS telegram_users (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    telegram_id BIGINT UNIQUE NOT NULL,
    username VARCHAR(100),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_telegram_users_user_id ON telegram_users(user_id);
CREATE INDEX IF NOT EXISTS idx_telegram_users_telegram_id ON telegram_users(telegram_id);
`,
}

// RunMigrations выполняет все миграции базы данных
func RunMigrations(db *sql.DB) error {
	for i, migration := range Migrations {
		log.Printf("Выполняется миграция %d...", i+1)
		_, err := db.Exec(migration)
		if err != nil {
			return fmt.Errorf("ошибка при выполнении миграции %d: %v", i+1, err)
		}
	}

	log.Println("Все миграции выполнены успешно")
	return nil
}

// CheckDatabaseConnection проверяет соединение с базой данных
func CheckDatabaseConnection(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return fmt.Errorf("не удалось подключиться к базе данных: %v", err)
	}
	return nil
}
