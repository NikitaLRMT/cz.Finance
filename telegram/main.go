package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"cz.Finance/telegram/client"
	"cz.Finance/telegram/handlers"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

func main() {
	// Загружаем .env файл, если он существует
	if err := godotenv.Load(); err != nil {
		logrus.Warn("Файл .env не найден, используются системные переменные окружения")
	}

	// Получаем токен бота из переменных окружения
	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramToken == "" {
		// Токен отсутствует, выводим предупреждение и завершаем работу
		logrus.Fatal("TELEGRAM_BOT_TOKEN не указан в переменных окружения. Используйте файл .env или экспортируйте переменную.")
	}

	// Получаем URL API из переменных окружения
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		// Если URL не указан, используем значение по умолчанию
		apiURL = "http://localhost:8080/api"
		logrus.Warnf("API_URL не указан в переменных окружения, используется значение по умолчанию: %s", apiURL)
	}

	// Создаем HTTP клиент для API
	apiClient := client.NewAPIClient(apiURL)

	// Настраиваем бота
	pref := telebot.Settings{
		Token:  telegramToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		logrus.Fatalf("Ошибка создания бота: %v", err)
	}

	// Инициализируем обработчики бота
	botHandlers := handlers.NewBotHandlers(apiClient)

	// Регистрируем обработчики
	botHandlers.RegisterHandlers(bot)

	// Запускаем бота в отдельной горутине
	go func() {
		logrus.Infof("Telegram бот запущен. Используется API URL: %s", apiURL)
		bot.Start()
	}()

	// Канал для сигналов останова
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Ожидаем сигнала для грацилозного завершения
	<-stop

	// Грацилозно завершаем бота
	logrus.Info("Останавливаем Telegram бота...")
	bot.Stop()

	logrus.Info("Telegram бот остановлен")
}
