package configs

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Config представляет собой структуру конфигурации приложения
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Logger   *logrus.Logger
}

// ServerConfig содержит настройки HTTP-сервера
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig содержит настройки базы данных
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	URL      string
}

// JWTConfig содержит настройки для JWT-аутентификации
type JWTConfig struct {
	Secret    string
	ExpiresIn time.Duration
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() *Config {
	// Загрузка переменных окружения из .env файла, если он существует
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("Файл .env не найден, используются системные переменные окружения")
	}

	// Настройка логгера
	logger := logrus.New()
	logLevel := getEnv("LOG_LEVEL", "info")
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(level)
	}
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Настройки HTTP-сервера
	readTimeout, _ := strconv.Atoi(getEnv("SERVER_READ_TIMEOUT", "10"))
	writeTimeout, _ := strconv.Atoi(getEnv("SERVER_WRITE_TIMEOUT", "10"))
	idleTimeout, _ := strconv.Atoi(getEnv("SERVER_IDLE_TIMEOUT", "60"))

	serverConfig := ServerConfig{
		Port:         getEnv("SERVER_PORT", "8080"),
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		IdleTimeout:  time.Duration(idleTimeout) * time.Second,
	}

	// Настройки базы данных
	dbConfig := DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "finance_app"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Настройки JWT
	jwtExpiresIn, _ := strconv.Atoi(getEnv("JWT_EXPIRES_IN", "24"))
	jwtConfig := JWTConfig{
		Secret:    getEnv("JWT_SECRET", "your-secret-key"),
		ExpiresIn: time.Duration(jwtExpiresIn) * time.Hour,
	}

	// Логирование JWT конфигурации
	logger.Infof("JWT ExpiresIn: %v часов", jwtExpiresIn)
	logger.Infof("JWT Secret length: %d символов", len(jwtConfig.Secret))
	logger.Infof("JWT Secret (first 3 chars): %s...", jwtConfig.Secret[:3])

	return &Config{
		Server:   serverConfig,
		Database: dbConfig,
		JWT:      jwtConfig,
		Logger:   logger,
	}
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// LoadFromEnv загружает конфигурацию из переменных окружения
func LoadFromEnv() (*Config, error) {
	// Пытаемся загрузить из .env файла, если он существует
	if err := godotenv.Load(); err != nil {
		logrus.Warn("Файл .env не найден, используются системные переменные окружения")
	}

	// Настройки базы данных
	dbConfig := &DatabaseConfig{
		Host:     getEnv("DB_HOST", "postgres"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "finance"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Проверяем DATABASE_URL - он имеет приоритет над отдельными настройками
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL != "" {
		dbConfig.URL = dbURL
	}

	// Настройки HTTP-сервера
	readTimeout, _ := strconv.Atoi(getEnv("SERVER_READ_TIMEOUT", "10"))
	writeTimeout, _ := strconv.Atoi(getEnv("SERVER_WRITE_TIMEOUT", "10"))
	idleTimeout, _ := strconv.Atoi(getEnv("SERVER_IDLE_TIMEOUT", "60"))

	serverConfig := ServerConfig{
		Port:         getEnv("SERVER_PORT", "8080"),
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		IdleTimeout:  time.Duration(idleTimeout) * time.Second,
	}

	// Настройки JWT
	jwtExpiresIn, _ := strconv.Atoi(getEnv("JWT_EXPIRES_IN", "24"))
	jwtConfig := JWTConfig{
		Secret:    getEnv("JWT_SECRET", "your-secret-key"),
		ExpiresIn: time.Duration(jwtExpiresIn) * time.Hour,
	}

	// Логирование JWT конфигурации
	logger := logrus.New()
	logLevel := getEnv("LOG_LEVEL", "info")
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(level)
	}
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.Infof("JWT ExpiresIn: %v часов", jwtExpiresIn)
	logger.Infof("JWT Secret length: %d символов", len(jwtConfig.Secret))
	logger.Infof("JWT Secret (first 3 chars): %s...", jwtConfig.Secret[:3])

	return &Config{
		Server:   serverConfig,
		Database: *dbConfig,
		JWT:      jwtConfig,
		Logger:   logger,
	}, nil
}
