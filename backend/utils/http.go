package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"cz.Finance/backend/models"

	"github.com/gorilla/mux"
)

// Типы для ключей контекста
type ContextKey string

const (
	UserIDKey ContextKey = "user_id"
	EmailKey  ContextKey = "email"
)

// RespondWithJSON отправляет ответ с данными в формате JSON
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Ошибка при сериализации JSON", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// RespondWithError отправляет ответ с ошибкой в формате JSON
func RespondWithError(w http.ResponseWriter, code int, message string, error string) {
	RespondWithJSON(w, code, models.ErrorResponse{
		Status:  code,
		Message: message,
		Error:   error,
	})
}

// ParseJSON парсит JSON из тела запроса в структуру
func ParseJSON(r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

// GetIDParam извлекает параметр id из URL
func GetIDParam(r *http.Request) (int64, error) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		return 0, errors.New("неверный формат ID")
	}
	return id, nil
}

// GetParams возвращает все параметры маршрута из URL
func GetParams(r *http.Request) map[string]string {
	return mux.Vars(r)
}

// GetQueryParam извлекает параметр из строки запроса
func GetQueryParam(r *http.Request, name string) string {
	return r.URL.Query().Get(name)
}

// GetIntQueryParam извлекает числовой параметр из строки запроса
func GetIntQueryParam(r *http.Request, name string, defaultValue int) int {
	param := GetQueryParam(r, name)
	if param == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(param)
	if err != nil {
		return defaultValue
	}

	return value
}

// GetUserIDFromContext извлекает ID пользователя из контекста запроса
func GetUserIDFromContext(r *http.Request) (int64, error) {
	userID, ok := r.Context().Value(UserIDKey).(int64)
	if !ok {
		return 0, errors.New("пользователь не аутентифицирован")
	}
	return userID, nil
}
