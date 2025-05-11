package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cz.Finance/backend/configs"
	"cz.Finance/backend/database"
	"cz.Finance/backend/routes"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Загружаем конфигурацию
	config, err := configs.LoadFromEnv()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Подключаемся к базе данных
	db, err := database.ConnectDB(config)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Выполняем миграции
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Ошибка выполнения миграций: %v", err)
	}

	// Создаем роутер
	router := mux.NewRouter()

	// Регистрируем маршруты
	routes.RegisterRoutes(router, db, config)

	// Создаем HTTP-сервер
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Настраиваем CORS
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Разрешаем запросы с любых доменов
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Максимальное время кэширования pre-flight запросов в секундах
	})

	// Создаем HTTP-сервер с CORS middleware
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      corsMiddleware.Handler(router),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Infof("Сервер запущен на порту %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	// Канал для сигналов останова
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Ожидаем сигнала для грацилозного завершения
	<-stop

	// Создаем контекст с таймаутом для грацилозного завершения
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Грацилозно завершаем сервер
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при остановке сервера: %v", err)
	}

	log.Info("Сервер остановлен")
}

