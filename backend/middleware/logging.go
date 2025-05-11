package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// LoggingMiddleware добавляет логирование запросов
func LoggingMiddleware(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Оборачиваем ResponseWriter для получения статус-кода
			responseWriter := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Передаем запрос дальше
			next.ServeHTTP(responseWriter, r)

			// Логируем информацию о запросе
			duration := time.Since(start)
			logger.WithFields(logrus.Fields{
				"method":     r.Method,
				"path":       r.URL.Path,
				"status":     responseWriter.statusCode,
				"duration":   duration.String(),
				"user_agent": r.UserAgent(),
				"ip":         getIP(r),
			}).Info("Request processed")
		})
	}
}

// responseWriter оборачивает стандартный ResponseWriter для получения статус-кода
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader сохраняет статус-код
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// getIP извлекает IP-адрес из запроса
func getIP(r *http.Request) string {
	// Проверяем X-Forwarded-For, который может быть установлен прокси-сервером
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return ip
	}

	// Если X-Forwarded-For не установлен, используем RemoteAddr
	return r.RemoteAddr
}
